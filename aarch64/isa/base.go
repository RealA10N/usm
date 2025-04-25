package aarch64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type NonBranchingInstruction struct{}

func (NonBranchingInstruction) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleContinue: true}, core.ResultList{}
}
