package usm64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type nonBranchingInstruction struct {
	baseInstruction
	NotCriticalInstruction
}

func newNonBranchingInstruction(info *gen.InstructionInfo) nonBranchingInstruction {
	return nonBranchingInstruction{baseInstruction: newBaseInstruction(info)}
}

func (i *nonBranchingInstruction) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleContinue: true}, core.ResultList{}
}
