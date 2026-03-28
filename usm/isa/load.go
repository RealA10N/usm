package usmisa

import (
	"fmt"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

// Load reads the value stored in a variable into a register.
//
// Syntax: <type> <reg> = load <var>
// Example: $64 %x = load &local
type Load struct {
	// Control Flow
	gen.NonBranchingInstruction

	// Dead Code Elimination
	opt.NonCriticalInstruction
	opt.UsesArgumentsInstruction
	opt.DefinesTargetsInstruction
}

func NewLoad() gen.InstructionDefinition {
	return Load{}
}

func (Load) Operator(*gen.InstructionInfo) string {
	return "load"
}

func (Load) Validate(info *gen.InstructionInfo) core.ResultList {
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
		// First use — infer variable type from the load target.
		varArg.Variable.Type = targetType
		return core.ResultList{}
	}

	if !targetType.Equal(varArg.Variable.Type) {
		varType := varArg.Variable.Type
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Target type does not match variable type",
				Location: info.Declaration,
			},
			{
				Type:     core.HintResult,
				Message:  fmt.Sprintf("Target type is \"%s\"", targetType),
				Location: targetType.Declaration,
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
