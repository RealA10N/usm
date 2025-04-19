package aarch64codegen

import "alon.kr/x/usm/gen"

// FunctionCodegenContext contains information about the code generation
// context for a specific function.
type FunctionCodegenContext struct {
	FileCodegenContext *FileCodegenContext

	// The offset of each basic block in the function, relative to the function
	// entry point.
	BlockOffsets map[*gen.BasicBlockInfo]uint64
}
