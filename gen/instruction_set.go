package gen

import (
	"alon.kr/x/faststringmap"
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// TODO: add basic interface methods for instruction.
type Instruction interface{}

// A basic instruction definition. This defines the logic that converts the
// generic, architecture / instruction set independent instruction AST nodes
// into a format instruction which is part of a specific instruction set.
type InstructionDefinition interface {

	// Returns the a constant slice with all valid names of the instruction.
	// This is called ones to initialize internal data structures that will
	// then be used to quickly point to the instruction definition.
	Names() []string

	// Build an instruction from the provided targets and arguments.
	BuildInstruction(
		targets []*RegisterInfo,
		arguments []*ArgumentInfo,
	) (Instruction, core.ResultList)

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
		ctx *GenerationContext,
		targets []*TypeInfo,
		arguments []*TypeInfo,
	) ([]*TypeInfo, core.ResultList)
}

type InstructionSet struct {
	NameToDefinition faststringmap.Map[InstructionDefinition]
}

func NewInstructionSet(instDefs []InstructionDefinition) InstructionSet {
	// optimization: # of entries is at least # of instructions.
	entries := make([]faststringmap.MapEntry[InstructionDefinition], 0, len(instDefs))

	for _, instDef := range instDefs {
		for _, name := range instDef.Names() {
			entry := faststringmap.MapEntry[InstructionDefinition]{
				Key:   name,
				Value: instDef,
			}
			entries = append(entries, entry)
		}
	}

	return InstructionSet{faststringmap.NewMap(entries)}
}

// Get the instruction definition that corresponds to the instruction in the
// provided parsed node, or return an error if the instruction is not known.
func (s *InstructionSet) getInstructionDefinitionFromNode(
	ctx *GenerationContext,
	node parse.InstructionNode,
) (InstructionDefinition, core.Result) {
	name := string(node.Operator.Raw(ctx.SourceContext))
	instDef, ok := s.NameToDefinition.LookupString(name)

	if !ok {
		return nil, core.GenericResult{
			Type:     core.ErrorResult,
			Message:  "Unknown instruction name",
			Location: &node.Operator,
		}
		// TODO: add typo suggestions?
	}

	return instDef, nil
}

func NewRegisterTypeMismatchError(
	NewDeclaration core.UnmanagedSourceView,
	FirstDeclaration core.UnmanagedSourceView,
) core.Result {
	return core.GenericResult{
		Type:     core.ErrorResult,
		Message:  "Explicit register type does not match previous declaration",
		Location: &NewDeclaration,
		Next: core.GenericResult{
			Type:     core.HintResult,
			Message:  "Previous declaration here",
			Location: &FirstDeclaration,
		},
	}
}

func (s *InstructionSet) getTargetTypeFromTargetNode(
	ctx *GenerationContext,
	node parse.TargetNode,
) (typeInfo *TypeInfo, res core.Result) {

	// if an explicit type is provided to the target, get the type info.
	var explicitType *TypeInfo
	if node.Type != nil {
		explicitTypeName := string(node.Type.Identifier.Raw(ctx.SourceContext))
		explicitType = ctx.Types.GetType(explicitTypeName)
	}

	registerName := string(node.Register.Raw(ctx.SourceContext))
	registerInfo := ctx.Registers.GetRegister(registerName)

	if registerInfo != nil {
		// register is already previously defined
		if explicitType != nil {
			// ensure explicit type matches the previously declared one.
			if explicitType != registerInfo.Type {
				return nil, NewRegisterTypeMismatchError(
					node.View(),
					registerInfo.Declaration,
				)
			}
		}

		// all checks passed; return previously defined register type.
		return registerInfo.Type, nil

	} else {
		// this is the first appearance of the register; if the type is provided
		// explicitly, use it. otherwise, there is no way to know the type of
		// the target register at this.
		// the type and register will be finalized when the instruction is built,
		// and only then it is added to the register manager.
		return explicitType, nil
	}
}

func (s *InstructionSet) getTargetTypesFromInstructionNode(
	ctx *GenerationContext,
	node parse.InstructionNode,
) ([]*TypeInfo, core.ResultList) {
	targets := make([]*TypeInfo, len(node.Targets))
	results := core.ResultList{}

	for i, target := range node.Targets {
		typeInfo, result := s.getTargetTypeFromTargetNode(ctx, target)
		targets[i] = typeInfo
		if result != nil {
			results.Append(result)
		}
	}

	return targets, results
}

func (s *InstructionSet) getArgumentFromArgumentNode(
	ctx *GenerationContext,
	node parse.ArgumentNode,
) (*ArgumentInfo, core.Result) {
	if registerNode, ok := node.(parse.RegisterNode); ok {
		// TODO: duplicated code: make function that extracts register name from node.
		registerName := string(registerNode.Raw(ctx.SourceContext))
		registerInfo := ctx.Registers.GetRegister(registerName)

		if registerInfo == nil {
			v := node.View()
			return nil, core.GenericResult{
				Type:     core.ErrorResult,
				Message:  "Undefined register used as argument",
				Location: &v,
			}
		}

		argumentInfo := ArgumentInfo{
			Type: registerInfo.Type,
		}
		return &argumentInfo, nil
	}

	v := node.View()
	return nil, core.GenericResult{
		Type:     core.InternalErrorResult,
		Message:  "Unsupported argument type",
		Location: &v,
	}
}

func (s *InstructionSet) getArgumentsFromInstructionNode(
	ctx *GenerationContext,
	node parse.InstructionNode,
) (arguments []*ArgumentInfo, results core.ResultList) {
	arguments = make([]*ArgumentInfo, len(node.Arguments))
	for i, arg := range node.Arguments {
		argInfo, result := s.getArgumentFromArgumentNode(ctx, arg)
		arguments[i] = argInfo
		if result != nil {
			results.Append(result)
		}
	}
	return
}

func (s *InstructionSet) argumentsToArgumentTypes(arguments []*ArgumentInfo) []*TypeInfo {
	argumentTypes := make([]*TypeInfo, len(arguments))
	for i, arg := range arguments {
		argumentTypes[i] = arg.Type
	}
	return argumentTypes
}

func (s *InstructionSet) defineNewRegister(
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
		return nil, core.GenericResult{
			Type:     core.InternalErrorResult,
			Message:  "internal register type mismatch",
			Location: &nodeView,
		}
	}

	return registerInfo, result
}

func (s *InstructionSet) defineNewRegisters(
	ctx *GenerationContext,
	node parse.InstructionNode,
	targetTypes []*TypeInfo,
) ([]*RegisterInfo, core.Result) {
	if len(node.Targets) != len(targetTypes) {
		v := node.View()
		return nil, core.GenericResult{
			Type:     core.InternalErrorResult,
			Message:  "targets length mismatch",
			Location: &v,
		}
	}

	registers := make([]*RegisterInfo, len(node.Targets))
	for i, target := range node.Targets {
		registerInfo, result := s.defineNewRegister(ctx, target, targetTypes[i])
		if result != nil {
			return nil, result
		}

		if registerInfo == nil {
			v := target.View()
			return nil, core.GenericResult{
				Type:     core.InternalErrorResult,
				Message:  "unexpected nil register",
				Location: &v,
			}
		}

		registers[i] = registerInfo
	}

	return registers, nil
}

// Convert an instruction parsed node into an instruction that is in the
// instruction set.
func (s *InstructionSet) Build(
	ctx *GenerationContext,
	node parse.InstructionNode,
) (inst Instruction, results core.ResultList) {
	instDef, res := s.getInstructionDefinitionFromNode(ctx, node)
	if res != nil {
		results.Append(res)
	}

	targetTypes, targetsResults := s.getTargetTypesFromInstructionNode(ctx, node)
	results.Extend(&targetsResults)

	arguments, argumentsResults := s.getArgumentsFromInstructionNode(ctx, node)
	results.Extend(&argumentsResults)

	if !results.IsEmpty() {
		return nil, results
	}

	argumentTypes := s.argumentsToArgumentTypes(arguments)
	actualTargetTypes, results := instDef.InferTargetTypes(ctx, targetTypes, argumentTypes)
	// TODO: validate that the returned target types matches expected constraints.

	if !results.IsEmpty() {
		return nil, results
	}

	targets, result := s.defineNewRegisters(ctx, node, actualTargetTypes)
	if result != nil {
		return nil, list.FromSlice([]core.Result{result})
	}

	return instDef.BuildInstruction(targets, arguments)
}
