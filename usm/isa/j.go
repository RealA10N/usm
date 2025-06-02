package usmisa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type J struct {
	// Control Flow
	gen.BranchToLabelArguments

	// Dead Code Elimination
	opt.CriticalInstruction
	opt.UsesNothingInstruction
	opt.DefinesNothingInstruction
}

func NewJump() J {
	return J{}
}

func (J) Operator(*gen.InstructionInfo) string {
	return "j"
}

func (i J) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	curResults := gen.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	curResults = gen.AssertArgumentsExactly(info, 1)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	labelArg := info.Arguments[0]
	_, curResults = gen.ArgumentToLabel(labelArg)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
