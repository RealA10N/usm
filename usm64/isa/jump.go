package usm64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type JumpInstruction struct {
	baseInstruction
}

func (i *JumpInstruction) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	label := i.InstructionInfo.Arguments[0].(*gen.LabelArgumentInfo).Label
	return gen.StepInfo{PossibleBranches: []*gen.LabelInfo{label}}, core.ResultList{}
}

func (i *JumpInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) core.ResultList {
	labelArgument := i.InstructionInfo.Arguments[0].(*gen.LabelArgumentInfo)
	return ctx.JumpToLabel(labelArgument.Label)
}

func (i *JumpInstruction) String() string {
	return "JMP"
}

func NewJumpInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	return gen.BaseInstruction(&JumpInstruction{baseInstruction: baseInstruction{info}}), core.ResultList{}
}

func NewJumpInstructionDefinition() gen.InstructionDefinition {
	return &FixedInstructionDefinition{
		Targets:   0,
		Arguments: 1,
		Creator:   NewJumpInstruction,
	}
}
