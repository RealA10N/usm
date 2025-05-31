package usmisa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Ret struct {
	gen.ReturningInstruction
}

func NewRet() gen.InstructionDefinition {
	return Ret{}
}

func (i Ret) Operator(*gen.InstructionInfo) string {
	return "ret"
}

func (i Ret) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	curResults := gen.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	curResults = gen.AssertArgumentsExactly(info, 0)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
