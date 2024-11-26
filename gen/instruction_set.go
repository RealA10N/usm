package gen

import (
	"alon.kr/x/faststringmap"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// TODO: add basic interface methods for instruction.
type Instruction interface {
	HasSideEffects() bool
}

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
		targets []RegisterInfo,
		arguments []ArgumentInfo,
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
	InterTargetTypes(
		targets []*TypeInfo,
		arguments []TypeInfo,
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

func (s *InstructionSet) getInstructionArgumentFromArgumentNode(
	ctx *GenerationContext,
	node parse.ArgumentNode,
) (ArgumentInfo, core.Result) {
	return nil, nil // TODO: implement
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

func (s *InstructionSet) getInstructionArgumentsFromInstructionNode(
	ctx *GenerationContext,
	node parse.InstructionNode,
) (args []ArgumentInfo, results core.ResultList) {
	for _, arg := range node.Arguments {
		info, res := s.getInstructionArgumentFromArgumentNode(ctx, arg)
		if res == nil {
			args = append(args, info)
		} else {
			results.Append(res)
		}
	}
	return
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

	targets, targetsResults := s.getInstructionTargetsFromInstructionNode(ctx, node)
	results.Extend(&targetsResults)

	arguments, argumentsResults := s.getInstructionArgumentsFromInstructionNode(ctx, node)
	results.Extend(&argumentsResults)

	if results.IsEmpty() {
		return instDef.BuildInstruction(targets, arguments)
	} else {
		return nil, results
	}
}
