package aarch64isa

import (
	"alon.kr/x/aarch64codegen/immediates"
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/aarch64codegen/registers"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Movz struct {
	gen.NonBranchingInstruction
}

func (Movz) Operator(*gen.InstructionInfo) string {
	return "movz"
}

func (i Movz) Xd(info *gen.InstructionInfo) (registers.GPRegister, core.ResultList) {
	results := aarch64translation.AssertTargetsExactly(info, 1)

	if !results.IsEmpty() {
		return registers.GPRegister(0), results
	}

	Xd, results := aarch64translation.TargetToAarch64GPRegister(info.Targets[0])
	if !results.IsEmpty() {
		return registers.GPRegister(0), results
	}

	return Xd, core.ResultList{}
}

func (i Movz) Immediate(
	info *gen.InstructionInfo,
) (immediates.Immediate16, instructions.MovShift, core.ResultList) {
	results := aarch64translation.AssertArgumentsBetween(info, 1, 2)
	if !results.IsEmpty() {
		return immediates.Immediate16(0), instructions.MovShift0, results

	}

	imm, curResults := aarch64translation.ArgumentToAarch64Immediate16(info.Arguments[0])
	results.Extend(&curResults)

	shift := instructions.MovShift0
	if len(info.Arguments) > 1 {
		shift, curResults = aarch64translation.ArgumentToAarch64MovShift(
			info.Arguments[1],
		)
		results.Extend(&curResults)
	}

	if !results.IsEmpty() {
		return immediates.Immediate16(0), instructions.MovShift0, results
	}

	return imm, shift, core.ResultList{}
}

func (i Movz) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	_, curResults := i.Xd(info)
	results.Extend(&curResults)

	_, _, curResults = i.Immediate(info)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}

func (i Movz) Codegen(
	_ *aarch64codegen.InstructionCodegenContext,
	info *gen.InstructionInfo,
) (instructions.Instruction, core.ResultList) {
	results := core.ResultList{}

	Xd, curResults := i.Xd(info)
	results.Extend(&curResults)

	imm, shift, curResults := i.Immediate(info)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	return instructions.MOVZ(Xd, imm, shift), core.ResultList{}
}
