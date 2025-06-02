package usmisa

import (
	"alon.kr/x/usm/gen"
)

type Xor struct {
	// Inharits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewXor() gen.InstructionDefinition {
	return Xor{}
}

func (Xor) Operator(*gen.InstructionInfo) string {
	return "xor"
}
