package usmisa

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// The call instruction
type Call struct {
	gen.NonBranchingInstruction
}

func NewCall() gen.InstructionDefinition {
	return Call{}
}

func (i Call) Operator(*gen.InstructionInfo) string {
	return "call"
}

func (i Call) IsCritical(*gen.InstructionInfo) bool {
	return true
}

func (i Call) Validate(info *gen.InstructionInfo) core.ResultList {
	results := gen.AssertAtLeastArguments(info, 1)
	if !results.IsEmpty() {
		return results
	}

	funcArg, ok := info.Arguments[0].(*gen.GlobalArgumentInfo)
	if !ok {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Call instruction requires a global function argument as the first argument",
				Location: funcArg.Declaration(),
			},
		})
	}

	funcInfo := info.FileInfo.GetFunction(funcArg.Name())
	if funcInfo == nil {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Call global function does not exist, or is not a function",
				Location: funcArg.Declaration(),
			},
		})
	}

	// TODO: this should also check types.
	curResults := gen.AssertArgumentsExactly(info, len(funcInfo.Parameters)+1)
	results.Extend(&curResults)

	curResults = gen.AssertTargetsExactly(info, len(funcInfo.Targets))
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
