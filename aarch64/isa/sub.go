package aarch64isa

import (
	"alon.kr/x/aarch64codegen/immediates"
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type BaseSub struct {
	NonBranchingInstruction
}

func (BaseSub) Operator() string {
	return "sub"
}

type Sub struct {
	BaseSub
	instructions.Sub
}

func (i Sub) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type SubImm struct {
	BaseSub
	instructions.SubImm
}

func (i SubImm) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type SubDefinition struct {
	immediates.SetFlags
}

func (d SubDefinition) buildRegisterVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := core.ResultList{}

	Xd, curResults := aarch64translation.TargetToAarch64GPRegister(info.Targets[0])
	results.Extend(&results)

	Xn, curResults := aarch64translation.ArgumentToAarch64GPRegister(info.Arguments[0])
	results.Extend(&curResults)

	// TODO: Add shifted register support
	Xm, results := aarch64translation.ArgumentToAarch64GPRegister(info.Arguments[1])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	return Sub{
		Sub: instructions.SUB(Xd, Xn, Xm, d.SetFlags),
	}, core.ResultList{}
}

func (SubDefinition) buildImmediateVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {

	results := core.ResultList{}

	Xd, curResults := aarch64translation.TargetToAarch64GPRegister(info.Targets[0])
	results.Extend(&results)

	Xn, curResults := aarch64translation.ArgumentToAarch64GPorSPRegister(info.Arguments[0])
	results.Extend(&curResults)

	// TODO: Add shifted immediate support
	imm, curResults := aarch64translation.ArgumentToAarch64Immediate12(info.Arguments[1])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	return SubImm{
		SubImm: instructions.SUBI(Xd, Xn, imm),
	}, core.ResultList{}
}

func (d SubDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := core.ResultList{}

	curResults := aarch64translation.AssertTargetsExactly(info, 1)
	results.Extend(&curResults)

	curResults = aarch64translation.AssertArgumentsExactly(info, 2)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	switch info.Arguments[1].(type) {
	case *gen.RegisterArgumentInfo:
		return d.buildRegisterVariant(info)

	case *gen.ImmediateInfo:
		return d.buildImmediateVariant(info)

	default:
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Second \"sub\" argument must be a register or immediate",
				Location: info.Arguments[1].Declaration(),
			},
		})
	}
}

func NewSubInstructionDefinition(setFlags immediates.SetFlags) gen.InstructionDefinition {
	return SubDefinition{SetFlags: setFlags}
}
