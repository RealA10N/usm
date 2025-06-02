package usmisa

import (
	"fmt"

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

	functionTargetTypes := info.FunctionInfo.Targets
	retTypes, curResults := gen.ArgumentsToTypes(info.Arguments)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	if len(functionTargetTypes) != len(retTypes) {
		results.Append(core.Result{
			{
				Type: core.ErrorResult,
				Message: fmt.Sprintf(
					"Number of arguments (%d) and the number of targets of the function (%d) do not match",
					len(retTypes),
					len(functionTargetTypes),
				),
				Location: info.Declaration,
			},
		})
	}

	if !results.IsEmpty() {
		return results
	}

	for i, funcTargetType := range functionTargetTypes {
		retArgumentType := retTypes[i]

		if !funcTargetType.Equal(retArgumentType) {
			results.Append(core.Result{
				{
					Type:     core.ErrorResult,
					Message:  "Return type does not match the type of the function target",
					Location: retArgumentType.Declaration,
				},
				{
					Type:     core.HintResult,
					Message:  "Matches this function target type",
					Location: funcTargetType.Declaration,
				},
			})
		}
	}

	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
