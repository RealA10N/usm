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

// isPointerTo reports whether targetType is exactly one pointer level deeper
// than varType.
//
// Because the descriptor *N packs N consecutive pointer levels into a single
// entry, adding one level of indirection to a type that already ends with *N
// produces *{N+1}, not a separate *1 appended.  The two cases are therefore:
//
//   - varType ends with *N  →  targetType is identical except the last
//     descriptor has Amount == N+1.
//   - varType does not end with *  →  targetType has the same descriptors as
//     varType plus a trailing *1.
func isPointerTo(targetType, varType gen.ReferencedTypeInfo) bool {
	if targetType.Base != varType.Base {
		return false
	}

	nVar := len(varType.Descriptors)
	nTarget := len(targetType.Descriptors)

	if nVar > 0 && varType.Descriptors[nVar-1].Type == gen.PointerTypeDescriptor {
		// varType ends with *N: targetType must be identical but with *(N+1).
		if nTarget != nVar {
			return false
		}
		for i := 0; i < nVar-1; i++ {
			if !targetType.Descriptors[i].Equal(varType.Descriptors[i]) {
				return false
			}
		}
		varLast := varType.Descriptors[nVar-1]
		targetLast := targetType.Descriptors[nTarget-1]
		expected := new(big.Int).Add(varLast.Amount, big.NewInt(1))
		return targetLast.Type == gen.PointerTypeDescriptor &&
			targetLast.Amount.Cmp(expected) == 0
	}

	// varType does not end with *: targetType must append a new *1.
	if nTarget != nVar+1 {
		return false
	}
	for i := 0; i < nVar; i++ {
		if !targetType.Descriptors[i].Equal(varType.Descriptors[i]) {
			return false
		}
	}
	last := targetType.Descriptors[nTarget-1]
	return last.Type == gen.PointerTypeDescriptor &&
		last.Amount.Cmp(big.NewInt(1)) == 0
}

// stripPointer returns t with one pointer level removed from the last
// descriptor.  If the last descriptor has Amount == 1 the descriptor is
// dropped entirely; otherwise its Amount is decremented by one.
// The caller must ensure t ends with a pointer descriptor.
func stripPointer(t gen.ReferencedTypeInfo) gen.ReferencedTypeInfo {
	n := len(t.Descriptors)
	last := t.Descriptors[n-1]

	newDescriptors := make([]gen.TypeDescriptorInfo, n)
	copy(newDescriptors, t.Descriptors)

	if last.Amount.Cmp(big.NewInt(1)) == 0 {
		newDescriptors = newDescriptors[:n-1]
	} else {
		newDescriptors[n-1] = gen.TypeDescriptorInfo{
			Type:   gen.PointerTypeDescriptor,
			Amount: new(big.Int).Sub(last.Amount, big.NewInt(1)),
		}
	}

	return gen.ReferencedTypeInfo{
		Base:        t.Base,
		Descriptors: newDescriptors,
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
		// First use — infer variable type by stripping one pointer level from
		// the lea target type.
		if len(targetType.Descriptors) == 0 ||
			targetType.Descriptors[len(targetType.Descriptors)-1].Type != gen.PointerTypeDescriptor {
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
