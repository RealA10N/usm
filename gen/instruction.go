package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// TODO: add basic interface methods for instruction.
type BaseInstruction interface{}

// A basic instruction definition. This defines the logic that converts the
// generic, architecture / instruction set independent instruction AST nodes
// into a format instruction which is part of a specific instruction set.
type InstructionDefinition[InstT BaseInstruction] interface {
	// Build an instruction from the provided targets and arguments.
	BuildInstruction(
		targets []*RegisterArgumentInfo,
		arguments []ArgumentInfo,
	) (InstT, core.ResultList)

	// Provided a list a list of types that correspond to argument types,
	// and a (possibly partial) list of target types, return a complete list
	// of target types which is implicitly inferred from the argument types,
	// and possibly the explicit target types, or an error if the target types
	// can not be inferred.
	//
	// On success, the length of the returned type slice should be equal to the
	// provided (partial) targets length. The non nil provided target types
	// should not be modified.
	//
	// TODO: perhaps we should not pass the bare generation context to the "public"
	// instruction set definition API, and should wrap it with a limited interface.
	InferTargetTypes(
		ctx *FunctionGenerationContext[InstT],
		targets []*ReferencedTypeInfo,
		arguments []*ReferencedTypeInfo,
	) ([]ReferencedTypeInfo, core.ResultList)
}

type InstructionInfo[InstT BaseInstruction] struct {
	// The actual instruction instance, which is part of the instruction set,
	// and not part of this package.
	Instruction InstT

	// The targets of the instruction.
	Targets []*RegisterArgumentInfo

	// The arguments of the instruction.
	Arguments []ArgumentInfo
}

// MARK: Manager

type InstructionManager[InstT BaseInstruction] interface {
	// Get the instruction definition that corresponds to the provided name.
	GetInstructionDefinition(name string) (InstructionDefinition[InstT], core.ResultList)
}

// MARK: Generator

type InstructionGenerator[InstT BaseInstruction] struct {
	ArgumentGenerator FunctionContextGenerator[InstT, parse.ArgumentNode, ArgumentInfo]
	TargetGenerator   FunctionContextGenerator[InstT, parse.TargetNode, registerPartialInfo]
}

func NewInstructionGenerator[InstT BaseInstruction]() FunctionContextGenerator[
	InstT,
	parse.InstructionNode,
	*InstructionInfo[InstT],
] {
	return FunctionContextGenerator[
		InstT,
		parse.InstructionNode,
		*InstructionInfo[InstT],
	](
		&InstructionGenerator[InstT]{
			ArgumentGenerator: NewArgumentGenerator[InstT](),
			TargetGenerator:   NewTargetGenerator[InstT](),
		},
	)
}

func (g *InstructionGenerator[InstT]) generateArguments(
	ctx *FunctionGenerationContext[InstT],
	node parse.InstructionNode,
) ([]ArgumentInfo, core.ResultList) {
	arguments := make([]ArgumentInfo, len(node.Arguments))
	results := core.ResultList{}

	// Different arguments should not effect one another.
	// Thus, we just collect all of the errors along the way, and return
	// them in one chunk.
	for i, argument := range node.Arguments {
		argInfo, curResults := g.ArgumentGenerator.Generate(ctx, argument)
		results.Extend(&curResults)
		arguments[i] = argInfo
	}

	return arguments, results
}

func (g *InstructionGenerator[InstT]) generatePartialTargetsInfo(
	ctx *FunctionGenerationContext[InstT],
	node parse.InstructionNode,
) ([]registerPartialInfo, core.ResultList) {
	targets := make([]registerPartialInfo, len(node.Targets))
	results := core.ResultList{}

	// Different targets should not effect one another.
	// Thus, we just collect all of the errors along the way, and return
	// them in one chunk.
	for i, target := range node.Targets {
		typeInfo, curResults := g.TargetGenerator.Generate(ctx, target)
		results.Extend(&curResults)
		targets[i] = typeInfo
	}

	return targets, results
}

func partialTargetsToTypes(targets []registerPartialInfo) []*ReferencedTypeInfo {
	types := make([]*ReferencedTypeInfo, len(targets))
	for i, target := range targets {
		types[i] = target.Type
	}
	return types
}

func argumentsToTypes(arguments []ArgumentInfo) []*ReferencedTypeInfo {
	types := make([]*ReferencedTypeInfo, len(arguments))
	for i, arg := range arguments {
		types[i] = arg.GetType()
	}
	return types
}

func (g *InstructionGenerator[InstT]) getTargetRegister(
	ctx *FunctionGenerationContext[InstT],
	node parse.TargetNode,
	targetType ReferencedTypeInfo,
) (*RegisterArgumentInfo, core.Result) {
	registerName := nodeToSourceString(ctx.FileGenerationContext, node.Register)
	registerInfo := ctx.Registers.GetRegister(registerName)
	nodeView := node.View()

	if registerInfo == nil {
		// register is defined here; we should create the register and define
		// it's type.
		newRegisterInfo := &RegisterInfo{
			Name:        registerName,
			Type:        targetType,
			Declaration: nodeView,
		}

		registerArgument := &RegisterArgumentInfo{
			Register:    newRegisterInfo,
			declaration: nodeView,
		}

		return registerArgument, ctx.Registers.NewRegister(newRegisterInfo)
	}

	// register is already defined
	if !registerInfo.Type.Equal(targetType) {
		// notest: sanity check only
		return nil, core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Internal register type mismatch",
			Location: &nodeView,
		}}
	}

	registerArgument := &RegisterArgumentInfo{
		Register:    registerInfo,
		declaration: nodeView,
	}

	return registerArgument, nil
}

// Registers can be defined by being a target of an instruction.
// After we have determined 100% of the instruction targets types (either
// if they were explicitly declared or not), we call this function with the
// target types, and here we iterate over all target types and define missing
// registers.
//
// This also returns the full list of register targets for the provided
// instruction.
func (g *InstructionGenerator[InstT]) defineAndGetTargetRegisters(
	ctx *FunctionGenerationContext[InstT],
	node parse.InstructionNode,
	targetTypes []ReferencedTypeInfo,
) ([]*RegisterArgumentInfo, core.ResultList) {
	if len(node.Targets) != len(targetTypes) {
		// notest: sanity check: ensure lengths match.
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Targets length mismatch",
			Location: &v,
		}})
	}

	registers := make([]*RegisterArgumentInfo, len(node.Targets))
	results := core.ResultList{}
	for i, target := range node.Targets {
		// register errors should not effect one another, so we collect them.
		registerInfo, result := g.getTargetRegister(ctx, target, targetTypes[i])
		if result != nil {
			results.Append(result)
		}

		if registerInfo == nil {
			// notest: sanity check, should not happen.
			v := target.View()
			results.Append(core.Result{{
				Type:     core.InternalErrorResult,
				Message:  "Unexpected nil register",
				Location: &v,
			}})
		}

		registers[i] = registerInfo
	}

	return registers, results
}

// Convert an instruction parsed node into an instruction that is in the
// instruction set.
// If new registers are defined in the instruction (by assigning values to
// instruction targets), the register is created and added to the generation
// context.
func (g *InstructionGenerator[InstT]) Generate(
	ctx *FunctionGenerationContext[InstT],
	node parse.InstructionNode,
) (info *InstructionInfo[InstT], results core.ResultList) {
	// We start generating the instruction, by getting the definition interface,
	// and processing the targets and arguments. We accumulate the results,
	// since those processes do not effect each other.

	instName := viewToSourceString(ctx.FileGenerationContext, node.Operator)
	instDef, results := ctx.Instructions.GetInstructionDefinition(instName)

	arguments, curResults := g.generateArguments(ctx, node)
	results.Extend(&curResults)

	partialTargets, curResults := g.generatePartialTargetsInfo(ctx, node)
	results.Extend(&curResults)

	// Now it's time to check if we have any errors so far.
	if !results.IsEmpty() {
		return
	}

	explicitTargetTypes := partialTargetsToTypes(partialTargets)
	argumentTypes := argumentsToTypes(arguments)
	targetTypes, results := instDef.InferTargetTypes(ctx, explicitTargetTypes, argumentTypes)
	// TODO: validate that the returned target types matches expected constraints.

	if !results.IsEmpty() {
		return
	}

	targets, results := g.defineAndGetTargetRegisters(ctx, node, targetTypes)
	if !results.IsEmpty() {
		return
	}

	instruction, results := instDef.BuildInstruction(targets, arguments)
	if !results.IsEmpty() {
		return
	}

	info = &InstructionInfo[InstT]{
		Instruction: instruction,
		Targets:     targets,
		Arguments:   arguments,
	}

	return info, core.ResultList{}
}
