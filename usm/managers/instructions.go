package usmmanagers

import (
	"alon.kr/x/faststringmap"
	"alon.kr/x/usm/gen"
	usmisa "alon.kr/x/usm/usm/isa"
)

func NewInstructionManager() gen.InstructionManager {
	return gen.NewInstructionMap(
		[]faststringmap.MapEntry[gen.InstructionDefinition]{
			// Arithmetic
			{Key: "", Value: usmisa.NewMove()},
			{Key: "add", Value: usmisa.NewAdd()},
			{Key: "sub", Value: usmisa.NewSub()},
			{Key: "mul", Value: usmisa.NewMul()},

			// Bitwise Operations
			{Key: "and", Value: usmisa.NewAnd()},
			{Key: "or", Value: usmisa.NewOr()},
			{Key: "xor", Value: usmisa.NewXor()},

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
