package aarch64codegen

import "alon.kr/x/usm/gen"

type BasicBlockCodegenContext struct {
	*FunctionCodegenContext
	*gen.BasicBlockInfo
}

func (ctx *BasicBlockCodegenContext) Offset() uint64 {
	return ctx.BasicBlockOffsets[ctx.BasicBlockInfo]
}
