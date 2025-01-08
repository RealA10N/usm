package usm64isa

import (
	"fmt"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type TerminateInstruction struct{}

func (i *TerminateInstruction) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleReturn: true}, core.ResultList{}
}

func (i *TerminateInstruction) String() string {
	return "TERM"
}

func (i *TerminateInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) core.ResultList {
	fmt.Println("[Terminate]")
	ctx.ShouldTerminate = true
	return core.ResultList{}
}

func NewTerminateInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	return gen.BaseInstruction(&TerminateInstruction{}), core.ResultList{}
}

func NewTerminateInstructionDefinition() gen.InstructionDefinition {
	return &FixedInstructionDefinition{
		Targets:   0,
		Arguments: 0,
		Creator:   NewTerminateInstruction,
	}
}
