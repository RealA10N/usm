package usmisa

import (
	"fmt"
	"math/big"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

// Lea (load effective address) computes the address of a variable and stores
// it in a pointer register. This bridges variables to the general pointer model:
// the resulting pointer can be passed to functions or stored in memory.
//
// Syntax: <type> * <reg> = lea <var>
// Example: $64 * %ptr = lea &local
type Lea struct {
	// Control Flow
	gen.NonBranchingInstruction

	// Dead Code Elimination
	opt.NonCriticalInstruction
	opt.UsesArgumentsInstruction
	opt.DefinesTargetsInstruction
}

func NewLea() gen.InstructionDefinition {
	return Lea{}
}

func (Lea) Operator(*gen.InstructionInfo) string {
	return "lea"
}

// isPointerTo reports whether targetType is exactly varType with one additional
// PointerTypeDescriptor appended (i.e. targetType == varType + one level of *).
func isPointerTo(targetType, varType gen.ReferencedTypeInfo) bool {
	if targetType.Base != varType.Base {
		return false
	}

	if len(targetType.Descriptors) != len(varType.Descriptors)+1 {
		return false
	}

	last := targetType.Descriptors[len(targetType.Descriptors)-1]
	if last.Type != gen.PointerTypeDescriptor || last.Amount.Cmp(big.NewInt(1)) != 0 {
		return false
	}

	for i, d := range varType.Descriptors {
		if targetType.Descriptors[i] != d {
			return false
		}
	}

	return true
}

func (Lea) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	curResults := gen.AssertTargetsExactly(info, 1)
	results.Extend(&curResults)

	curResults = gen.AssertArgumentsExactly(info, 1)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	varArg, curResults := gen.ArgumentToVariable(info.Arguments[0])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	targetType := info.Targets[0].Register.Type
	varType := varArg.Variable.Type

	if !isPointerTo(targetType, varType) {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Target type must be a pointer to the variable type",
				Location: info.Declaration,
			},
			{
				Type:     core.HintResult,
				Message:  fmt.Sprintf("Target type is \"%s\"", targetType),
				Location: targetType.Declaration,
			},
			{
				Type:     core.HintResult,
				Message:  fmt.Sprintf("Expected \"%s *\"", varType),
				Location: varType.Declaration,
			},
		})
	}

	return core.ResultList{}
}
