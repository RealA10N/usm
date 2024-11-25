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
	Names() []string
	Builder(targets []RegisterInfo, arguments []ArgumentInfo) (Instruction, core.ResultList)
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
	ctx core.SourceContext,
	node parse.InstructionNode,
) (InstructionDefinition, core.Result) {
	name := string(node.Operator.Raw(ctx))
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
	ctx core.SourceContext,
	node parse.TargetNode,
) (RegisterInfo, core.Result) {
	return RegisterInfo{}, nil // TODO: implement
}

func (s *InstructionSet) getInstructionArgumentFromArgumentNode(
	ctx core.SourceContext,
	node parse.ArgumentNode,
) (ArgumentInfo, core.Result) {
	return nil, nil // TODO: implement
}

func (s *InstructionSet) getInstructionTargetsFromInstructionNode(
	ctx core.SourceContext,
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
	ctx core.SourceContext,
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
	ctx core.SourceContext,
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
		return instDef.Builder(targets, arguments)
	} else {
		return nil, results
	}
}
