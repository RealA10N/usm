package usmisa

import (
	"alon.kr/x/usm/gen"
)

// Jump to label if the value is not negative (positive or zero).
type Jnn struct {
	// Inharits most of the functionality from ConditionalJump
	ConditionalJump
}

func NewJnn() gen.InstructionDefinition {
	return Jnn{}
}

func (Jnn) Operator(*gen.InstructionInfo) string {
	return "jnn"
}
