package usm64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type JumpZeroInstruction struct {
	baseInstruction
}

func (i *JumpZeroInstruction) PossibleNextSteps() ([]gen.StepInfo, core.ResultList) {
	label := i.InstructionInfo.Arguments[1].(*gen.LabelArgumentInfo).Label
	return []gen.StepInfo{
		gen.ContinueToNextInstruction{},
		gen.JumpToLabel{Label: label},
	}, core.ResultList{}
}

func (i *JumpZeroInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) core.ResultList {
	value, results := ctx.ArgumentToValue(i.InstructionInfo.Arguments[0])
	if !results.IsEmpty() {
		return results
	}

	if value == uint64(0) {
		labelArgument := i.InstructionInfo.Arguments[1].(*gen.LabelArgumentInfo)
		return ctx.JumpToLabel(labelArgument.Label)

	} else {
		return ctx.ContinueToNextInstruction()
	}
}

func NewJumpZeroInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	return gen.BaseInstruction(&JumpZeroInstruction{
		baseInstruction: newBaseInstruction(info),
	}), core.ResultList{}
}

func NewJumpZeroInstructionDefinition() gen.InstructionDefinition {
	return &FixedInstructionDefinition{
		Targets:   0,
		Arguments: 2,
		Creator:   NewJumpZeroInstruction,
	}
}
