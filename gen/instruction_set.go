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
	Builder(targets []RegisterInfo, arguments []ArgumentInfo) (Instruction, core.Result)
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

func (s *InstructionSet) getInstructionTargetsFromNode(
	ctx core.SourceContext,
	node parse.InstructionNode,
) []RegisterInfo {
	return nil // TODO: implement
}

func (s *InstructionSet) getInstructionArgumentsFromNode(
	ctx core.SourceContext,
	node parse.InstructionNode,
) []ArgumentInfo {
	return nil // TODO: implement
}

// Convert an instruction parsed node into an instruction that is in the
// instruction set.
func (s *InstructionSet) Build(
	ctx core.SourceContext,
	node parse.InstructionNode,
) (Instruction, core.Result) {
	instDef, err := s.getInstructionDefinitionFromNode(ctx, node)
	if err != nil {
		return nil, err
	}

	targets := s.getInstructionTargetsFromNode(ctx, node)
	arguments := s.getInstructionArgumentsFromNode(ctx, node)
	return instDef.Builder(targets, arguments)
}
