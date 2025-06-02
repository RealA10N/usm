package usmisa

import (
	"alon.kr/x/usm/gen"
)

type Or struct {
	// Inharits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewOr() gen.InstructionDefinition {
	return Or{}
}

func (Or) Operator(*gen.InstructionInfo) string {
	return "or"
}
