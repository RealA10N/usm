package usm64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

// Currently, supporting only a single argument + target.
// TODO: support n arguments and target pairs.
type MovInstruction struct {
	Target usm64core.Register
	Value  usm64core.ValuedArgument
}

func (i *MovInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) usm64core.EmulationError {
	ctx.Registers[i.Target.Name] = i.Value.Value(ctx)
	ctx.IncrementInstructionPointer()
	return nil
}

func NewMovInstruction(
	targets []usm64core.Register,
	arguments []usm64core.Argument,
) (usm64core.Instruction, core.ResultList) {
	results := core.ResultList{}

	value, valueResults := usm64core.ArgumentToValuedArgument(arguments[0])
	results.Extend(&valueResults)

	if !results.IsEmpty() {
		return nil, results
	}

	return &MovInstruction{Target: targets[0], Value: value}, core.ResultList{}
}

func NewMovInstructionDefinition() gen.InstructionDefinition[usm64core.Instruction] {
	return &FixedInstructionDefinition{
		Targets:   1,
		Arguments: 1,
		Creator:   NewMovInstruction,
	}
}
