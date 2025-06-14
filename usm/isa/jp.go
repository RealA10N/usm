package usmisa

import (
	"alon.kr/x/usm/gen"
)

// Jump to label if the value is positive.
type Jp struct {
	// Inherits most of the functionality from ConditionalJump
	ConditionalJump
}

func NewJp() gen.InstructionDefinition {
	return Jp{}
}

func (Jp) Operator(*gen.InstructionInfo) string {
	return "jp"
}
