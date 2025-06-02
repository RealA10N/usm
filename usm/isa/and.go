package usmisa

import (
	"alon.kr/x/usm/gen"
)

type And struct {
	// Inharits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewAnd() gen.InstructionDefinition {
	return And{}
}

func (And) Operator(*gen.InstructionInfo) string {
	return "and"
}
