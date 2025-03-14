package usm64isa

import (
	"fmt"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type PutInstruction struct {
	nonBranchingInstruction
	CriticalInstruction
}

func (i *PutInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) core.ResultList {
	value, results := ctx.ArgumentToValue(i.Arguments[0])
	if !results.IsEmpty() {
		return results
	}

	fmt.Println(value)
	return ctx.ContinueToNextInstruction()
}

func (i *PutInstruction) String() string {
	return "PUT"
}

func NewPutInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	return gen.BaseInstruction(&PutInstruction{
		nonBranchingInstruction: newNonBranchingInstruction(info),
	}), core.ResultList{}
}

func NewPutInstructionDefinition() gen.InstructionDefinition {
	return &FixedInstructionDefinition{
		Targets:   0,
		Arguments: 1,
		Creator:   NewPutInstruction,
	}
}
