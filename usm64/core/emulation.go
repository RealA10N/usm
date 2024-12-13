package usm64core

// MARK: Error

type EmulationError interface{}

// MARK: Context

type EmulationContext struct {
	// The next instruction index to execute.
	// Should be len(instructions) to indicate the return from the function.
	NextInstructionIndex uint64

	Registers map[string]uint64
}

func NewEmulationContext() EmulationContext {
	return EmulationContext{
		NextInstructionIndex: 0,
		Registers:            make(map[string]uint64),
	}
}

type Emulateable interface {
	Emulate(ctx *EmulationContext) EmulationError
}
