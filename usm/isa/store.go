package usmisa

import (
	"fmt"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

// Store writes the value of a register into a variable.
//
// Syntax: store <var> <reg>
// Example: store &local %x
type Store struct {
	// Control Flow
	gen.NonBranchingInstruction

	// Dead Code Elimination: store is critical because it has a side effect
	// (writes to the local frame) that must not be eliminated.
	opt.CriticalInstruction
	opt.UsesArgumentsInstruction
	opt.DefinesNothingInstruction
}

func NewStore() gen.InstructionDefinition {
	return Store{}
}

func (Store) Operator(*gen.InstructionInfo) string {
	return "store"
}

func (Store) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	curResults := gen.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	curResults = gen.AssertArgumentsExactly(info, 2)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	varArg, curResults := gen.ArgumentToVariable(info.Arguments[0])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	valueType, curResults := gen.ArgumentToType(info.Arguments[1])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	varType := varArg.Variable.Type

	if !valueType.Equal(varType) {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Value type does not match variable type",
				Location: info.Declaration,
			},
			{
				Type:     core.HintResult,
				Message:  fmt.Sprintf("Value type is \"%s\"", valueType),
				Location: valueType.Declaration,
			},
			{
				Type:     core.HintResult,
				Message:  fmt.Sprintf("Variable type is \"%s\"", varType),
				Location: varType.Declaration,
			},
		})
	}

	return core.ResultList{}
}
