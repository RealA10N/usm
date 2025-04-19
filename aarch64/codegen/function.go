package aarch64codegen

import "alon.kr/x/usm/gen"

// FunctionCodegenContext contains information about the code generation
// context for a specific function.
type FunctionCodegenContext struct {
	FileCodegenContext *FileCodegenContext

	// The offset of each basic block in the function, relative to the function
	// entry point.
	BasicBlockOffsets map[*gen.BasicBlockInfo]uint64
}

func NewFunctionCodegenContext(
	fileCtx *FileCodegenContext,
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
		FileCodegenContext: fileCtx,
		BasicBlockOffsets:  basicBlockOffsets,
	}
}
