package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Add struct {
	gen.NonBranchingInstruction
}

func NewAdd() gen.InstructionDefinition {
	return Add{}
}

func (Add) Operator(*gen.InstructionInfo) string {
	return "add"
}

func (add Add) codegenRegisterVariant(
	info *gen.InstructionInfo,
) (instructions.Instruction, core.ResultList) {
	Xd, Xn, Xm, results := aarch64translation.BinaryInstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	inst := instructions.NewAddShiftedRegister(Xd, Xn, Xm)
	return inst, core.ResultList{}
}

func (add Add) codegenImmediateVariant(
	info *gen.InstructionInfo,
) (instructions.Instruction, core.ResultList) {
	Xd, Xn, imm, results := aarch64translation.Immediate12InstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	inst := instructions.NewAddImmediate(Xd, Xn, imm)
	return inst, core.ResultList{}
}

func (add Add) Codegen(
	ctx *aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	// TODO: this implementation is very similar to the one in adds.go, and possibly
	// other binary arithmetic instructions. Consider refactoring this.

	info := ctx.InstructionInfo
	results := aarch64translation.ValidateBinaryInstruction(info)
	if !results.IsEmpty() {
		return nil, results
	}

	switch info.Arguments[1].(type) {
	case *gen.RegisterArgumentInfo:
		return add.codegenRegisterVariant(info)
	case *gen.ImmediateInfo:
		return add.codegenImmediateVariant(info)
	default:
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Second \"add\" argument must be a register or immediate",
				Location: info.Arguments[1].Declaration(),
			},
		})
	}
}

func (add Add) Validate(
	info *gen.InstructionInfo,
) core.ResultList {
	// TODO: this is a pretty hacky way to validate the instruction: we create
	// a "mock" generation context, and then try to generate the binary
	// representation of the instruction.
	ctx := aarch64codegen.InstructionCodegenContext{InstructionInfo: info}
	_, results := add.Codegen(&ctx)
	return results
}
