package aarch64codegen

import (
	"bytes"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// FileCodegenContext contains information about the code generation
// context for a specific file.
type FileCodegenContext struct {
	*gen.FileInfo

	// The offset of each function in the file (object file), relative to the
	// base object offset.
	FunctionOffsets map[*gen.FunctionInfo]uint64
}

func NewFileCodegenContext(file *gen.FileInfo) *FileCodegenContext {
	functionOffsets := make(map[*gen.FunctionInfo]uint64, len(file.Functions))

	offset := uint64(0)
	for _, function := range file.Functions {
		functionOffsets[function] = offset
		functionSize := uint64(function.Size()) * 4 // TODO: handle overflow?
		offset += functionSize
	}

	return &FileCodegenContext{
		FileInfo:        file,
		FunctionOffsets: functionOffsets,
	}
}

func (ctx *FileCodegenContext) newFunctionCodegenContext(
	function *gen.FunctionInfo,
) *FunctionCodegenContext {
	basicBlocks := function.CollectBasicBlocks()
	basicBlockOffsets := make(map[*gen.BasicBlockInfo]uint64, len(basicBlocks))

	offset := uint64(0)
	for _, block := range basicBlocks {
		basicBlockOffsets[block] = offset

		// In AArch64, each instruction is of constant size of 4 bytes.
		// TODO: handle overflow?
		offset += uint64(block.Size()) * 4
	}

	return &FunctionCodegenContext{
		FileCodegenContext: ctx,
		FunctionInfo:       function,
		BasicBlockOffsets:  basicBlockOffsets,
	}
}

func (ctx *FileCodegenContext) Codegen(
	buffer *bytes.Buffer,
) core.ResultList {
	for _, function := range ctx.Functions {
		funcCtx := ctx.newFunctionCodegenContext(function)
		results := funcCtx.Codegen(buffer)
		if !results.IsEmpty() {
			return results
		}
	}

	return core.ResultList{}
}
