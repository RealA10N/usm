package gen

type InstructionInfo[InstT BaseInstruction] struct {
	// The actual instruction instance, which is part of the instruction set,
	// and not part of this package.
	Instruction InstT

	// The targets of the instruction.
	Targets []*RegisterArgumentInfo

	// The arguments of the instruction.
	Arguments []ArgumentInfo
}
