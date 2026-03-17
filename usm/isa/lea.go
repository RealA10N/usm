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
		if !targetType.Descriptors[i].Equal(d) {
			return false
		}
	}

	return true
}

// stripPointer returns t with the last (pointer) descriptor removed.
// The caller must ensure t is a pointer type.
func stripPointer(t gen.ReferencedTypeInfo) gen.ReferencedTypeInfo {
	return gen.ReferencedTypeInfo{
		Base:        t.Base,
		Descriptors: t.Descriptors[:len(t.Descriptors)-1],
		Declaration: t.Declaration,
	}
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

	if varArg.Variable.Type.Base == nil {
		// First use — infer variable type by stripping the pointer from the
		// lea target type.
		last := targetType.Descriptors[len(targetType.Descriptors)-1]
		if len(targetType.Descriptors) == 0 ||
			last.Type != gen.PointerTypeDescriptor ||
			last.Amount.Cmp(big.NewInt(1)) != 0 {
			return list.FromSingle(core.Result{
				{
					Type:     core.ErrorResult,
					Message:  "lea target must be a pointer type",
					Location: info.Declaration,
				},
				{
					Type:     core.HintResult,
					Message:  fmt.Sprintf("Target type is \"%s\"", targetType),
					Location: targetType.Declaration,
				},
			})
		}
		varArg.Variable.Type = stripPointer(targetType)
		return core.ResultList{}
	}

	if !isPointerTo(targetType, varArg.Variable.Type) {
		varType := varArg.Variable.Type
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
