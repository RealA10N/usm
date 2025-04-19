package aarch64codegen

import "alon.kr/x/usm/gen"

type InstructionCodegenContext struct {
	*BasicBlockCodegenContext
	*gen.InstructionInfo

	InstructionOffsetInBasicBlock uint64
}

func Offset(ctx *InstructionCodegenContext) uint64 {
	basicBlockOffset := ctx.BasicBlockCodegenContext.Offset()
	return basicBlockOffset + ctx.InstructionOffsetInBasicBlock
}
