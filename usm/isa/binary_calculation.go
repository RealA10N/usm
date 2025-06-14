package usmisa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type BinaryCalculation struct {
	// Control Flow
	gen.NonBranchingInstruction

	// Dead Code Elimination
	opt.NonCriticalInstruction
	opt.UsesArgumentsInstruction
	opt.DefinesTargetsInstruction
}

// Validates that there are exactly two arguments and one target, all of the
// same type.
func (BinaryCalculation) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	curResults := gen.AssertTargetsExactly(info, 1)
	results.Extend(&curResults)

	curResults = gen.AssertArgumentsExactly(info, 2)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	leftType, curResults := gen.ArgumentToType(info.Arguments[0])
	results.Extend(&curResults)

	rightType, curResults := gen.ArgumentToType(info.Arguments[1])
	results.Extend(&curResults)

	targetType := gen.TargetToType(info.Targets[0])

	if !results.IsEmpty() {
		return results
	}

	if !leftType.Equal(targetType) {
		results.Append(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Left argument type does not match target type",
				Location: info.Arguments[0].Declaration(),
			},
		})
	}

	if !rightType.Equal(targetType) {
		results.Append(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Right argument type does not match target type",
				Location: info.Arguments[1].Declaration(),
			},
		})
	}

	return results
}
