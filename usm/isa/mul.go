package usmisa

import (
	"alon.kr/x/usm/gen"
)

type Mul struct {
	// Inherits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewMul() gen.InstructionDefinition {
	return Mul{}
}

func (Mul) Operator(*gen.InstructionInfo) string {
	return "mul"
}
