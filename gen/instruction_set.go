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

func (s *InstructionSet) getInstructionTargetFromTargetNode(
	ctx *GenerationContext,
	node parse.TargetNode,
) (regInfo RegisterInfo, res core.Result) {
	// hintedType := getTargetTypeFromTargetNode(ctx, node)
	if node.Type != nil {
		typeName := string(node.Type.Identifier.Raw(ctx.SourceContext))
		regInfo.Type = ctx.Types.GetType(typeName)
		if regInfo.Type == nil {
			res = core.GenericResult{
				Type:     core.ErrorResult,
				Message:  "Undefined type",
				Location: &node.Type.Identifier,
			}
			return
		}
	}

	registerName := string(node.Register.Raw(ctx.SourceContext))
	if existingReg := ctx.Registers.GetRegister(registerName); existingReg != nil {
		if regInfo.Type != nil {

		}
	}
	return
}

func (s *InstructionSet) getInstructionArgumentFromArgumentNode(
	ctx *GenerationContext,
	node parse.ArgumentNode,
) (ArgumentInfo, core.Result) {
	return nil, nil // TODO: implement
}

func (s *InstructionSet) getInstructionTargetsFromInstructionNode(
	ctx *GenerationContext,
	node parse.InstructionNode,
) (regs []RegisterInfo, results core.ResultList) {
	for _, target := range node.Targets {
		info, res := s.getInstructionTargetFromTargetNode(ctx, target)
		if res == nil {
			regs = append(regs, info)
		} else {
			results.Append(res)
		}
	}
	return
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
