package usmisa

import (
	"alon.kr/x/usm/gen"
)

type Sub struct {
	// Inherits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewSub() gen.InstructionDefinition {
	return Sub{}
}

func (Sub) Operator(*gen.InstructionInfo) string {
	return "sub"
}
