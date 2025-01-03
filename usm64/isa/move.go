package usm64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

// Currently, supporting only a single argument + target.
// TODO: support n arguments and target pairs.
type MoveInstruction struct {
	nonBranchingInstruction
}

func (i *MoveInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) core.ResultList {
	value, results := ctx.ArgumentToValue(i.Arguments[0])
	if !results.IsEmpty() {
		return results
	}

	targetName := i.Targets[0].Register.Name
	ctx.Registers[targetName] = value
	return ctx.ContinueToNextInstruction()
}

func (i *MoveInstruction) String() string {
	return ""
}

func NewMoveInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	return gen.BaseInstruction(&MoveInstruction{
		nonBranchingInstruction: newNonBranchingInstruction(info),
	}), core.ResultList{}
}

func NewMoveInstructionDefinition() gen.InstructionDefinition {
	return &FixedInstructionDefinition{
		Targets:   1,
		Arguments: 1,
		Creator:   NewMoveInstruction,
	}
}
