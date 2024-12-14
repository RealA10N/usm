package usm64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type AddInstruction struct {
	Target        usm64core.Register
	First, Second usm64core.ValuedArgument
}

func (i *AddInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) usm64core.EmulationError {
	ctx.Registers[i.Target.Name] = i.First.Value(ctx) + i.Second.Value(ctx)
	ctx.IncrementInstructionPointer()
	return nil
}

func NewAddInstruction(
	targets []usm64core.Register,
	arguments []usm64core.Argument,
) (usm64core.Instruction, core.ResultList) {
	results := core.ResultList{}

	first, firstResults := usm64core.ArgumentToValuedArgument(arguments[0])
	results.Extend(&firstResults)

	second, secondResults := usm64core.ArgumentToValuedArgument(arguments[1])
	results.Extend(&secondResults)

	if !results.IsEmpty() {
		return nil, results
	}

	return &AddInstruction{
		Target: targets[0],
		First:  first,
		Second: second,
	}, core.ResultList{}
}

func NewAddInstructionDefinition() gen.InstructionDefinition[usm64core.Instruction] {
	return &FixedInstructionDefinition{
		Targets:   1,
		Arguments: 2,
		Creator:   NewAddInstruction,
	}
}
