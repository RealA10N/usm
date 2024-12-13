package usm64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type AddInstruction struct {
	Target        usm64core.Register
	First, Second usm64core.Argument
}

func (i *AddInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) usm64core.EmulationError {
	ctx.Registers[i.Target] = i.First.Value(ctx) + i.Second.Value(ctx)
	ctx.NextInstructionIndex++
	return nil
}

func NewAddInstruction(
	targets []usm64core.Register,
	arguments []usm64core.Argument,
) (usm64core.Instruction, core.ResultList) {
	return &AddInstruction{
		Target: targets[0],
		First:  arguments[0],
		Second: arguments[1],
	}, core.ResultList{}
}

func NewAddInstructionDefinition() gen.InstructionDefinition[usm64core.Instruction] {
	return &FixedInstructionDefinition{
		Targets:   1,
		Arguments: 2,
		Creator:   NewAddInstruction,
	}
}
