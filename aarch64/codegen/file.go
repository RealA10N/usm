package aarch64codegen

import "alon.kr/x/usm/gen"

// FileCodegenContext contains information about the code generation
// context for a specific file.
type FileCodegenContext struct {
	// The offset of each function in the file (object file), relative to the
	// base object offset.
	FunctionOffsets map[*gen.FunctionInfo]uint64
}

func NewFileCodegenContext(file *gen.FileInfo) *FileCodegenContext {
	functionOffsets := make(map[*gen.FunctionInfo]uint64, len(file.Functions))

	offset := uint64(0)
	for _, function := range file.Functions {
		functionOffsets[function] = offset
		functionSize := uint64(function.Size() * 4) // TODO: handle overflow?
		offset += functionSize
	}

	return &FileCodegenContext{
		FunctionOffsets: functionOffsets,
	}
}
