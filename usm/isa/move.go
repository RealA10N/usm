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
	opt.UsesArgumentsInstruction
	opt.DefinesTargetsInstruction
}

// PropagateConstants returns the immediate constant assigned to the target
// register if the source argument is a constant immediate, otherwise nil.
func (Move) PropagateConstants(info *gen.InstructionInfo) []opt.ConstantDefinition {
	if len(info.Arguments) != 1 || len(info.Targets) != 1 {
		return nil
	}

	immediate, ok := info.Arguments[0].(*gen.ImmediateInfo)
	if !ok {
		return nil
	}

	return []opt.ConstantDefinition{{
		Register:  info.Targets[0].Register,
		Immediate: immediate,
	}}
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

	targetType := info.Targets[0].Register.Type

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

func (Move) Defines(info *gen.InstructionInfo) []*gen.RegisterInfo {
	return gen.TargetsToRegisters(info.Targets)
}

func (Move) Uses(info *gen.InstructionInfo) []*gen.RegisterInfo {
	return gen.ArgumentsToRegisters(info.Arguments)
}
