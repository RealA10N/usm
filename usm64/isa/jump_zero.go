package usm64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type JumpZeroInstruction struct {
	Argument usm64core.ValuedArgument
	Label    usm64core.Label
}

func (i *JumpZeroInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) usm64core.EmulationError {
	if i.Argument.Value(ctx) == uint64(0) {
		ctx.JumpToLabel(i.Label)
	} else {
		ctx.IncrementInstructionPointer()
	}
	return nil
}

func NewJumpZeroInstruction(
	targets []usm64core.Register,
	arguments []usm64core.Argument,
) (usm64core.Instruction, core.ResultList) {
	results := core.ResultList{}

	argument, argumentResults := usm64core.ArgumentToValuedArgument(arguments[0])
	results.Extend(&argumentResults)

	label, labelResults := usm64core.ArgumentToLabel(arguments[1])
	results.Extend(&labelResults)

	if !results.IsEmpty() {
		return nil, results
	}

	return &JumpZeroInstruction{Argument: argument, Label: label}, core.ResultList{}
}

func NewJumpZeroInstructionDefinition() gen.InstructionDefinition[usm64core.Instruction] {
	return &FixedInstructionDefinition{
		Targets:   0,
		Arguments: 2,
		Creator:   NewJumpZeroInstruction,
	}
}
