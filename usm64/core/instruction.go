package usm64core

// MARK: Instruction

type Instruction interface {
	// Emulate (interpret) the instruction, provided the context.
	Emulateable
}
