package aarch64codegen

import (
	"bytes"

	"alon.kr/x/macho/load/section64"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// FileCodegenContext contains information about the code generation
// context for a specific file.
type FileCodegenContext struct {
	*gen.FileInfo

	// The offset of each function in the file (object file), relative to the
	// base object offset. It stores offsets of only defined functions.
	FunctionOffsets map[*gen.FunctionInfo]uint64

	// The index of each function, according to their order in the source
	// definition. Indices start from 0 and are continuous.
	FunctionIndices map[*gen.FunctionInfo]uint32

	// The list of static relocations needed to be applied to the produced
	// object file before it is linked to an executable.
	Relocations []section64.RelocationBuilder
}

func NewFileCodegenContext(file *gen.FileInfo) *FileCodegenContext {
	functionOffsets := make(map[*gen.FunctionInfo]uint64)
	functionIndices := make(map[*gen.FunctionInfo]uint32, len(file.Functions))

	offset := uint64(0)
	idx := uint32(0)
	for _, function := range file.Functions {
		if function.IsDefined() {
			functionOffsets[function] = offset
			functionIndices[function] = idx

			functionSize := uint64(function.Size()) * 4 // TODO: handle overflow?
			offset += functionSize
		}

		idx++
	}

	return &FileCodegenContext{
		FileInfo:        file,
		FunctionOffsets: functionOffsets,
		FunctionIndices: functionIndices,
		Relocations:     []section64.RelocationBuilder{},
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
