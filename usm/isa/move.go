package usmisa

import (
	"fmt"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type Move struct {
	// Control Flow
	gen.NonBranchingInstruction

	// Dead Code Elimination
	opt.NonCriticalInstruction
	opt.UsesInstruction
	opt.DefinesTargetsInstruction
}

func NewMove() gen.InstructionDefinition {
	return Move{}
}

func (Move) Operator(*gen.InstructionInfo) string {
	return ""
}

func (Move) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	curResults := gen.AssertTargetsExactly(info, 1)
	results.Extend(&curResults)

	curResults = gen.AssertArgumentsExactly(info, 1)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	argumentType, curResults := gen.ArgumentToType(info.Arguments[0])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	targetType, curResults := gen.ArgumentToType(info.Targets[0])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	if !targetType.Equal(argumentType) {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Target and argument types do not match",
				Location: info.Declaration,
			},
			{
				Type:     core.HintResult,
				Message:  fmt.Sprintf("Target type is \"%s\"", targetType),
				Location: targetType.Declaration,
			},
			{
				Type:     core.HintResult,
				Message:  fmt.Sprintf("Argument type is \"%s\"", argumentType),
				Location: argumentType.Declaration,
			},
		})
	}

	return core.ResultList{}
}
