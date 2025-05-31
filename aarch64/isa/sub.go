package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Sub struct {
	gen.NonBranchingInstruction
}

func NewSubInstruction() gen.InstructionDefinition {
	return Sub{}
}

func (Sub) Operator(*gen.InstructionInfo) string {
	return "sub"
}

func (Sub) codegenRegisterVariant(
	info *gen.InstructionInfo,
) (instructions.Instruction, core.ResultList) {
	Xd, Xn, Xm, results := aarch64translation.BinaryInstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	inst := instructions.NewSubShiftedRegister(Xd, Xn, Xm)
	return inst, core.ResultList{}
}

func (Sub) codegenImmediateVariant(
	info *gen.InstructionInfo,
) (instructions.Instruction, core.ResultList) {
	Xd, Xn, imm, results := aarch64translation.Immediate12GPRegisterTargetInstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	inst := instructions.NewSubImmediate(Xd, Xn, imm)
	return inst, core.ResultList{}
}

func (i Sub) Codegen(
	ctx *aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	info := ctx.InstructionInfo
	results := aarch64translation.ValidateBinaryInstruction(info)
	if !results.IsEmpty() {
		return nil, results
	}

	switch info.Arguments[1].(type) {
	case *gen.RegisterArgumentInfo:
		return i.codegenRegisterVariant(info)
	case *gen.ImmediateInfo:
		return i.codegenImmediateVariant(info)
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

func (i Sub) Validate(
	info *gen.InstructionInfo,
) core.ResultList {
	// TODO: this is a pretty hacky way to validate the instruction: we create
	// a "mock" generation context, and then try to generate the binary
	// representation of the instruction.
	ctx := aarch64codegen.InstructionCodegenContext{InstructionInfo: info}
	_, results := i.Codegen(&ctx)
	return results
}
