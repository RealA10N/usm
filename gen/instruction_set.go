package gen

import (
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
		targets []*RegisterInfo,
		arguments []*ArgumentInfo,
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
		ctx *GenerationContext[InstT],
		targets []*TypeInfo,
		arguments []*TypeInfo,
	) ([]*TypeInfo, core.ResultList)
}

// MARK: Manager

type InstructionManager[InstT BaseInstruction] interface {
	// Get the instruction definition that corresponds to the provided name.
	GetInstructionDefinition(name string) (InstructionDefinition[InstT], core.ResultList)
}

// MARK: Generator

type InstructionGenerator[InstT BaseInstruction] struct {
	ArgumentGenerator ArgumentGenerator[InstT]
	TargetGenerator   TargetGenerator[InstT]
}

func (g *InstructionGenerator[InstT]) generateArguments(
	ctx *GenerationContext[InstT],
	node parse.InstructionNode,
) ([]*ArgumentInfo, core.ResultList) {
	arguments := make([]*ArgumentInfo, len(node.Arguments))
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

func (g *InstructionGenerator[InstT]) generateExplicitTargetsTypes(
	ctx *GenerationContext[InstT],
	node parse.InstructionNode,
) ([]*TypeInfo, core.ResultList) {
	targets := make([]*TypeInfo, len(node.Targets))
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

func argumentsToArgumentTypes(arguments []*ArgumentInfo) []*TypeInfo {
	types := make([]*TypeInfo, len(arguments))
	for i, arg := range arguments {
		types[i] = arg.Type
	}
	return types
}

func (g *InstructionGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.InstructionNode,
) (inst InstT, results core.ResultList) {
	// We start generating the instruction, by getting the definition interface,
	// and processing the targets and arguments. We accumulate the results,
	// since those processes do not effect each other.

	instName := string(node.Operator.Raw(ctx.SourceContext))
	instDef, results := ctx.Instructions.GetInstructionDefinition(instName)

	arguments, curResults := g.generateArguments(ctx, node)
	results.Extend(&curResults)

	explicitTargets, curResults := g.generateExplicitTargetsTypes(ctx, node)
	results.Extend(&curResults)

	// Now it's time to check if we have any errors so far.
	if !results.IsEmpty() {
		return
	}

	argumentTypes := argumentsToArgumentTypes(arguments)
	_, results = instDef.InferTargetTypes(ctx, explicitTargets, argumentTypes)
	// TODO: validate that the returned target types matches expected constraints.

	if !results.IsEmpty() {
		return
	}

	// return instDef.BuildInstruction(targets, arguments)
	return // TODO: finish implementing this
}

// MARK: old code

func (s *InstructionSet[InstT]) defineNewRegister(
	ctx *GenerationContext,
	node parse.TargetNode,
	targetType *TypeInfo,
) (registerInfo *RegisterInfo, result core.Result) {
	registerName := string(node.Register.Raw(ctx.SourceContext))
	registerInfo = ctx.Registers.GetRegister(registerName)
	nodeView := node.View()

	if registerInfo == nil {
		// register is defined here! we should create the register and define
		// it's type.
		newRegisterInfo := &RegisterInfo{
			Name:        registerName,
			Type:        targetType,
			Declaration: nodeView,
		}

		result = ctx.Registers.NewRegister(newRegisterInfo)
		return newRegisterInfo, result

	}

	// register is already defined;
	// sanity check: ensure the type matches the previously defined one.
	if registerInfo.Type != targetType {
		return nil, core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "internal register type mismatch",
			Location: &nodeView,
		}}
	}

	return registerInfo, result
}

func (s *InstructionSet[InstT]) defineNewRegisters(
	ctx *GenerationContext,
	node parse.InstructionNode,
	targetTypes []*TypeInfo,
) ([]*RegisterInfo, core.Result) {
	if len(node.Targets) != len(targetTypes) {
		v := node.View()
		return nil, core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "targets length mismatch",
			Location: &v,
		}}
	}

	registers := make([]*RegisterInfo, len(node.Targets))
	for i, target := range node.Targets {
		registerInfo, result := s.defineNewRegister(ctx, target, targetTypes[i])
		if result != nil {
			return nil, result
		}

		if registerInfo == nil {
			v := target.View()
			return nil, core.Result{{
				Type:     core.InternalErrorResult,
				Message:  "unexpected nil register",
				Location: &v,
			}}
		}

		registers[i] = registerInfo
	}

	return registers, nil
}
