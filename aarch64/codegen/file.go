package aarch64codegen

import "alon.kr/x/usm/gen"

// FileCodegenContext contains information about the code generation
// context for a specific file.
type FileCodegenContext struct {
	// The offset of each function in the file (object file), relative to the
	// base object offset.
	FunctionOffsets map[*gen.FunctionInfo]uint64
}
