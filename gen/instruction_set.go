package gen

import (
	"fmt"

	"alon.kr/x/faststringmap"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// TODO: add basic interface methods for instruction.
type Instruction interface {
	HasSideEffects() bool
}

type InstructionBuilder func(
	targets []parse.ParameterNode,
	arguments []parse.ArgumentNode,
) (Instruction, error)

// A basic instruction definition. This defines the logic that converts the
// generic, architecture / instruction set independent instruction AST nodes
// into a format instruction which is part of a specific instruction set.
type InstructionDef struct {
	Names   []string
	Builder InstructionBuilder
}

type InstructionSet struct {
	nameToBuilder faststringmap.Map[InstructionBuilder]
}

func NewInstructionSet(instructionDefs []InstructionDef) InstructionSet {
	// optimization: # of entries is at least # of instructions.
	entries := make([]faststringmap.MapEntry[InstructionBuilder], 0, len(instructionDefs))

	for _, inst := range instructionDefs {
		for _, name := range inst.Names {
			entry := faststringmap.MapEntry[InstructionBuilder]{
				Key:   name,
				Value: inst.Builder,
			}
			entries = append(entries, entry)
		}
	}

	return InstructionSet{faststringmap.NewMap(entries)}
}

// Convert an instruction parsed node into an instruction that is in the
// instruction set.
func (set *InstructionSet) Build(
	ctx core.SourceContext,
	node parse.InstructionNode,
) (Instruction, error) {
	var ok bool
	name := string(node.Operator.Raw(ctx))
	instBuilder, ok := set.nameToBuilder.LookupString(name)
	if !ok {
		return nil, fmt.Errorf("unknown instruction '%s'", name)
	}

	return instBuilder(node.Targets, node.Arguments)
}
