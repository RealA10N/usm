package usmisa

import (
	"fmt"

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
		// First use — infer variable type by dereferencing the lea target.
		varType, ok := targetType.Deref()
		if !ok {
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
		varArg.Variable.Type = varType
		return core.ResultList{}
	}

	if !varArg.Variable.Type.PointerTo().Equal(targetType) {
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
				Message:  fmt.Sprintf("Expected \"%s\"", varType.PointerTo()),
				Location: varType.Declaration,
			},
		})
	}

	return core.ResultList{}
}
