package usm64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type JumpNotZeroInstruction struct {
	baseInstruction
	CriticalInstruction
}

func (i *JumpNotZeroInstruction) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	label := i.InstructionInfo.Arguments[1].(*gen.LabelArgumentInfo).Label
	return gen.StepInfo{
		PossibleBranches: []*gen.LabelInfo{label},
		PossibleContinue: true,
	}, core.ResultList{}
}

func (i *JumpNotZeroInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) core.ResultList {
	value, results := ctx.ArgumentToValue(i.InstructionInfo.Arguments[0])
	if !results.IsEmpty() {
		return results
	}

	if value != uint64(0) {
		labelArgument := i.InstructionInfo.Arguments[1].(*gen.LabelArgumentInfo)
		return ctx.JumpToLabel(labelArgument.Label)

	} else {
		return ctx.ContinueToNextInstruction()
	}
}

func (i *JumpNotZeroInstruction) String() string {
	return "JNZ"
}

func NewJumpNotZeroInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	return gen.BaseInstruction(&JumpNotZeroInstruction{
		baseInstruction: newBaseInstruction(info),
	}), core.ResultList{}
}

func NewJumpNotZeroInstructionDefinition() gen.InstructionDefinition {
	return &FixedInstructionDefinition{
		Targets:   0,
		Arguments: 2,
		Creator:   NewJumpNotZeroInstruction,
	}
}
