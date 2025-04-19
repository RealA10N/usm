package aarch64isa

import (
	"math/big"

	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/usm/gen"
)

// FileCodegenContext contains information about the code generation
// context for a specific file.
type FileCodegenContext struct {
	// The offset of each function in the file (object file), relative to the
	// base object offset.
	FunctionOffsets map[*gen.FunctionInfo]*big.Int
}

// FunctionCodegenContext contains information about the code generation
// context for a specific function.
type FunctionCodegenContext struct {
	FileCodegenContext *FileCodegenContext

	// The offset of each basic block in the function, relative to the function
	// entry point.
	BlockOffsets map[*gen.BasicBlockInfo]*big.Int
}

type InstructionGenerator interface {
	// Converts the abstract instruction representation into a concrete binary
	// instruction.
	Generate(*FunctionCodegenContext) instructions.Instruction
}
