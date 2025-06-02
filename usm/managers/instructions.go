package usmmanagers

import (
	"alon.kr/x/faststringmap"
	"alon.kr/x/usm/gen"
	usmisa "alon.kr/x/usm/usm/isa"
)

func NewInstructionManager() gen.InstructionManager {
	return gen.NewInstructionMap(
		[]faststringmap.MapEntry[gen.InstructionDefinition]{
			{Key: "", Value: usmisa.NewMove()},

			// Functions
			{Key: "ret", Value: usmisa.NewRet()},
			{Key: "call", Value: usmisa.NewCall()},

			// Control Flow
			{Key: "j", Value: usmisa.NewJump()},  // Unconditional jump
			{Key: "jz", Value: usmisa.NewJz()},   // jump if zero
			{Key: "jnz", Value: usmisa.NewJnz()}, // jump if not zero
			{Key: "jp", Value: usmisa.NewJp()},   // jump if positive
			{Key: "jnp", Value: usmisa.NewJnp()}, // jump if not positive
			{Key: "jn", Value: usmisa.NewJp()},   // jump if negative
			{Key: "jnn", Value: usmisa.NewJnp()}, // jump if not negative

			// Static Single Assignment (SSA)
			{Key: "phi", Value: usmisa.NewPhi()},
		},
		false,
	)
}
