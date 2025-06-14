package usmisa

import (
	"alon.kr/x/usm/gen"
)

// Jump to label if the value is not positive (negative or zero).
type Jnp struct {
	// Inherits most of the functionality from ConditionalJump
	ConditionalJump
}

func NewJnp() gen.InstructionDefinition {
	return Jnp{}
}

func (Jnp) Operator(*gen.InstructionInfo) string {
	return "jnp"
}
