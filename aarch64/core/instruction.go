package aarch64core

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/usm/gen"
)

// FileCodegenContext contains information about the code generation
// context for a specific file.
type FileCodegenContext struct {
	// The offset of each function in the file (object file), relative to the
	// base object offset.
	FunctionOffsets map[*gen.FunctionInfo]uint64
}

// FunctionCodegenContext contains information about the code generation
// context for a specific function.
type FunctionCodegenContext struct {
	FileCodegenContext *FileCodegenContext

	// The offset of each basic block in the function, relative to the function
	// entry point.
	BlockOffsets map[*gen.BasicBlockInfo]uint64
}

type Instruction interface {
	gen.BaseInstruction

	// Converts the abstract instruction representation into a concrete binary
	// instruction.
	Generate(*FunctionCodegenContext) instructions.Instruction
}
