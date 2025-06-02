package usmisa

import (
	"alon.kr/x/usm/gen"
)

// Jump to label if the value is non-zero.
type Jnz struct {
	// Inharits most of the functionality from ConditionalJump
	ConditionalJump
}

func NewJnz() gen.InstructionDefinition {
	return Jnz{}
}

func (Jnz) Operator(*gen.InstructionInfo) string {
	return "jnz"
}
