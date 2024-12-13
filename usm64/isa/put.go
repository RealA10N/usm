package usm64isa

import (
	"fmt"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

// MARK: Debug

type PutInstruction struct {
	Argument usm64core.Argument
}

func (i *PutInstruction) Emulate(ctx *usm64core.EmulationContext) usm64core.EmulationError {
	fmt.Println(i.Argument.Value(ctx))
	ctx.NextInstructionIndex++
	return nil
}

func NewPutInstruction(
	targets []usm64core.Register,
	argument []usm64core.Argument,
) (usm64core.Instruction, core.ResultList) {
	return &PutInstruction{
		Argument: argument[0],
	}, core.ResultList{}
}

func NewPutInstructionDefinition() gen.InstructionDefinition[usm64core.Instruction] {
	return &FixedInstructionDefinition{
		Targets:   0,
		Arguments: 1,
		Creator:   NewPutInstruction,
	}
}
