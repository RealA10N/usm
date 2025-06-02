package usmisa

import (
	"alon.kr/x/usm/gen"
)

// Jump to label if the value is zero.
type Jz struct {
	// Inharits most of the functionality from ConditionalJump
	ConditionalJump
}

func NewJz() gen.InstructionDefinition {
	return Jz{}
}

func (Jz) Operator(*gen.InstructionInfo) string {
	return "jz"
}
