package usmisa

import (
	"alon.kr/x/list"
	"alon.kr/x/set"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type Phi struct {
	// Control Flow
	gen.NonBranchingInstruction

	// Dead Code Elimination
	opt.NonCriticalInstruction
	opt.UsesArgumentsInstruction
	opt.DefinesTargetsInstruction
}

func NewPhi() gen.InstructionDefinition {
	return Phi{}
}

func (Phi) Operator(*gen.InstructionInfo) string {
	return "phi"
}

func (Phi) validateEvenArguments(info *gen.InstructionInfo) core.ResultList {
	if len(info.Arguments)%2 != 0 {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "A \"phi\" instruction must have an even number of arguments.",
				Location: info.Declaration,
			},
			{
				Type:    core.HintResult,
				Message: "Each pair of arguments should consist of a label and a value.",
			},
		})
	}

	return core.ResultList{}
}

func (i Phi) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	curResults := gen.AssertTargetsExactly(info, 1)
	results.Extend(&curResults)

	curResults = i.validateEvenArguments(info)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	targetType := info.Targets[0].Register.Type
	incomingEdges := set.FromSlice(info.BasicBlockInfo.BackwardEdges)

	for i := 0; i < len(info.Arguments); i += 2 {
		labelArg := info.Arguments[i]

		label, curResults := gen.ArgumentToLabel(labelArg)
		results.Extend(&curResults)

		valueArg := info.Arguments[i+1]
		argType, curResults := gen.ArgumentToType(valueArg)
		results.Extend(&curResults)

		if !results.IsEmpty() {
			continue
		}

		block := label.BasicBlock
		if incomingEdges.Contains(block) {
			incomingEdges.Remove(block)
		} else {
			results.Append(core.Result{
				{
					Type:     core.ErrorResult,
					Message:  "Label does not match any incoming edge, or appears multiple times.",
					Location: labelArg.Declaration(),
				},
			})
		}

		if !argType.Equal(targetType) {
			results.Append(core.Result{
				{
					Type:     core.ErrorResult,
					Message:  "Argument type does not match target type.",
					Location: valueArg.Declaration(),
				},
			})
		}
	}

	// Notice that we do not check if incomingEdges is empty.
	// This is because it it VALID for some incoming edges to not have a
	// specified value. In that case the value will be undefined.

	return results
}

func (Phi) AddForwardingRegister(
	instruction *gen.InstructionInfo,
	block *gen.BasicBlockInfo,
	register *gen.RegisterInfo,
) core.ResultList {
	labelArg := gen.NewLabelArgumentInfo(block.Label)
	regArg := gen.NewRegisterArgumentInfo(register)
	instruction.AppendArgument(labelArg, regArg)
	return core.ResultList{}
}
