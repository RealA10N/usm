package usmisa

import (
	"math/big"

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

	// Constant Propagation: default is no folding; specific operations
	// override PropagateConstants to implement constant folding.
	opt.PropagatesNoConstants
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

// foldBinaryConstants evaluates a binary instruction when both arguments are
// immediates, returning the folded constant for the single target.
// fold receives the left and right big.Int values and returns the result.
// Returns nil if either argument is not an immediate.
func foldBinaryConstants(
	info *gen.InstructionInfo,
	fold func(l, r *big.Int) *big.Int,
) []opt.ConstantDefinition {
	if len(info.Arguments) != 2 || len(info.Targets) != 1 {
		return nil
	}

	left, ok := info.Arguments[0].(*gen.ImmediateInfo)
	if !ok {
		return nil
	}

	right, ok := info.Arguments[1].(*gen.ImmediateInfo)
	if !ok {
		return nil
	}

	result := fold(left.Value, right.Value)
	imm := &gen.ImmediateInfo{
		Type:  info.Targets[0].Register.Type,
		Value: result,
	}

	return []opt.ConstantDefinition{{
		Register:  info.Targets[0].Register,
		Immediate: imm,
	}}
}
