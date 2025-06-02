package usmisa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type Ret struct {
	// Control Flow
	gen.ReturningInstruction

	// Dead Code Elimination
	opt.CriticalInstruction
	opt.UsesArgumentsInstruction
	opt.DefinesNothingInstruction
}

func NewRet() gen.InstructionDefinition {
	return Ret{}
}

func (Ret) Operator(*gen.InstructionInfo) string {
	return "ret"
}

func (Ret) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	curResults := gen.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	// TODO: this is not exactly correct, arguments to the ret instruction must
	// match the targets of the function.
	curResults = gen.AssertArgumentsExactly(info, 0)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
