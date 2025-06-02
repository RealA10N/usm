package usmisa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type ConditionalJump struct {
	// Control Flow
	gen.BranchesToLabelArgumentsOrContinues

	// Dead Code Elimination
	opt.CriticalInstruction
	opt.UsesArgumentsInstruction
	opt.DefinesNothingInstruction
}

func (ConditionalJump) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	curResults := gen.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	curResults = gen.AssertArgumentsExactly(info, 2)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	valueArg := info.Arguments[0]
	_, curResults = gen.ArgumentToType(valueArg)
	results.Extend(&curResults)

	labelArg := info.Arguments[1]
	_, curResults = gen.ArgumentToLabel(labelArg)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
