package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Movz struct {
	NonBranchingInstruction
	instructions.Movz
}

func (Movz) Operator() string {
	return "movz"
}

func (i Movz) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type MovzDefinition struct{}

func (MovzDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := core.ResultList{}

	curResults := aarch64translation.AssertTargetsExactly(info, 1)
	results.Extend(&curResults)

	curResults = aarch64translation.AssertArgumentsBetween(info, 1, 2)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	Xd, curResults := aarch64translation.TargetToAarch64GPRegister(info.Targets[0])
	results.Extend(&curResults)

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
		return nil, results
	}

	return Movz{
		Movz: instructions.MOVZ(Xd, imm, shift),
	}, core.ResultList{}
}

func NewMovzInstructionDefinition() gen.InstructionDefinition {
	return MovzDefinition{}
}
