package usm64isa

import "alon.kr/x/usm/gen"

type baseInstruction struct {
	// A pointer the the internal USM representation of the instruction.
	// This in turn has the representation of the arguments, targets, and types.
	*gen.InstructionInfo
}

func newBaseInstruction(info *gen.InstructionInfo) baseInstruction {
	return baseInstruction{info}
}
