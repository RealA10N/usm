package usmisa

import (
	"alon.kr/x/usm/gen"
)

type Add struct {
	// Inherits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewAdd() gen.InstructionDefinition {
	return Add{}
}

func (Add) Operator(*gen.InstructionInfo) string {
	return "add"
}
