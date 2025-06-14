package usmisa

import (
	"alon.kr/x/usm/gen"
)

// Jump to label if the value is negative.
type Jn struct {
	// Inherits most of the functionality from ConditionalJump
	ConditionalJump
}

func NewJn() gen.InstructionDefinition {
	return Jn{}
}

func (Jn) Operator(*gen.InstructionInfo) string {
	return "jn"
}
