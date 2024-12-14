package usm64isa

import (
	"fmt"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type PutInstruction struct {
	Argument usm64core.ValuedArgument
}

func (i *PutInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) usm64core.EmulationError {
	fmt.Println(i.Argument.Value(ctx))
	ctx.NextInstructionIndex++
	return nil
}

func NewPutInstruction(
	targets []usm64core.Register,
	argument []usm64core.Argument,
) (usm64core.Instruction, core.ResultList) {
	valued, results := usm64core.ArgumentToValuedArgument(argument[0])
	if !results.IsEmpty() {
		return nil, results
	}

	return &PutInstruction{Argument: valued}, core.ResultList{}
}

func NewPutInstructionDefinition() gen.InstructionDefinition[usm64core.Instruction] {
	return &FixedInstructionDefinition{
		Targets:   0,
		Arguments: 1,
		Creator:   NewPutInstruction,
	}
}
