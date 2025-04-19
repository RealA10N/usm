package aarch64codegen

import (
	"bytes"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// FunctionCodegenContext contains information about the code generation
// context for a specific function.
type FunctionCodegenContext struct {
	*FileCodegenContext
	*gen.FunctionInfo

	// The offset of each basic block in the function, relative to the function
	// entry point.
	BasicBlockOffsets map[*gen.BasicBlockInfo]uint64
}

func (ctx *FunctionCodegenContext) newInstructionCodegenContext(
	instruction *gen.InstructionInfo,
	instructionOffsetInFunction uint64,
) *InstructionCodegenContext {
	return &InstructionCodegenContext{
		FunctionCodegenContext:      ctx,
		InstructionInfo:             instruction,
		InstructionOffsetInFunction: instructionOffsetInFunction,
	}
}

func (ctx *FunctionCodegenContext) Codegen(
	buffer *bytes.Buffer,
) core.ResultList {
	instructions := ctx.FunctionInfo.CollectInstructions()
	for idx, inst := range instructions {
		instOffset := uint64(idx * 4) // TODO: handle overflow?
		instCtx := ctx.newInstructionCodegenContext(
			inst,
			instOffset,
		)

		results := instCtx.Codegen(buffer)
		if !results.IsEmpty() {
			return results
		}
	}

	return core.ResultList{}
}
