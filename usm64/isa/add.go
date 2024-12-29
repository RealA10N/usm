package usm64isa

import (
	"math/bits"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type AddInstruction struct {
	nonBranchingInstruction
}

func (i *AddInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) core.ResultList {
	results := core.ResultList{}

	first, firstResults := ctx.ArgumentToValue(i.Arguments[0])
	results.Extend(&firstResults)

	second, secondResults := ctx.ArgumentToValue(i.Arguments[1])
	results.Extend(&secondResults)

	if !results.IsEmpty() {
		return results
	}

	targetName := i.Targets[0].Register.Name
	sum, _ := bits.Add64(first, second, 0)
	ctx.Registers[targetName] = sum
	return ctx.ContinueToNextInstruction()
}

func NewAddInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := core.ResultList{}

	if !results.IsEmpty() {
		return nil, results
	}

	return gen.BaseInstruction(&AddInstruction{
		nonBranchingInstruction: newNonBranchingInstruction(info),
	}), core.ResultList{}
}

func NewAddInstructionDefinition() gen.InstructionDefinition {
	return &FixedInstructionDefinition{
		Targets:   1,
		Arguments: 2,
		Creator:   NewAddInstruction,
	}
}
