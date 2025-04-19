package aarch64core

import (
	"alon.kr/x/aarch64codegen/instructions"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	"alon.kr/x/usm/gen"
)

type Instruction interface {
	gen.BaseInstruction

	// Converts the abstract instruction representation into a concrete binary
	// instruction.
	Generate(*aarch64codegen.FunctionCodegenContext) instructions.Instruction
}
