package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/aarch64codegen/registers"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Ret struct{}

func NewRet() gen.InstructionDefinition {
	return Ret{}
}

func (Ret) Operator(*gen.InstructionInfo) string {
	return "ret"
}

func (Ret) PossibleNextSteps(*gen.InstructionInfo) (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleReturn: true}, core.ResultList{}
}

func (i Ret) Codegen(
	ctx *aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	info := ctx.InstructionInfo
	Xn, results := i.Xn(info)
	if !results.IsEmpty() {
		return nil, results
	}

	return instructions.RET(Xn), core.ResultList{}
}

func (i Ret) Xn(info *gen.InstructionInfo) (registers.GPRegister, core.ResultList) {
	results := aarch64translation.AssertArgumentsBetween(info, 0, 1)
	if !results.IsEmpty() {
		return registers.GPRegister(0), results
	}

	Xn := registers.GPRegisterX30
	if len(info.Arguments) > 0 {
		Xn, results = aarch64translation.ArgumentToAarch64GPRegister(info.Arguments[0])
		if !results.IsEmpty() {
			return registers.GPRegister(0), results
		}
	}

	return Xn, core.ResultList{}
}

func (i Ret) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	_, curResults := i.Xn(info)
	results.Extend(&curResults)

	curResults = aarch64translation.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
