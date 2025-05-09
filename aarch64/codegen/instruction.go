package aarch64codegen

import (
	"bytes"
	"encoding/binary"

	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Instruction interface {
	gen.InstructionDefinition

	// Converts the abstract instruction representation into a concrete binary
	// instruction.
	Codegen(
		ctx *InstructionCodegenContext,
	) (instructions.Instruction, core.ResultList)
}

type InstructionCodegenContext struct {
	*FunctionCodegenContext
	*gen.InstructionInfo

	InstructionOffsetInFunction uint64
}

func (ctx *InstructionCodegenContext) InstructionOffsetInFile() uint64 {
	return ctx.InstructionOffsetInFunction + ctx.FunctionOffsets[ctx.FunctionInfo]
}

func (ctx *InstructionCodegenContext) Codegen(
	buffer *bytes.Buffer,
) core.ResultList {
	instruction, ok := ctx.InstructionInfo.Instruction.(Instruction)
	if !ok {
		return list.FromSingle(core.Result{
			{
				Type:     core.InternalErrorResult,
				Message:  "Instruction is not an AArch64 instruction",
				Location: ctx.Declaration,
			},
		})
	}

	binaryInst, results := instruction.Codegen(ctx)
	if !results.IsEmpty() {
		return results
	}

	binary.Write(buffer, binary.LittleEndian, binaryInst.Binary())
	return core.ResultList{}
}
