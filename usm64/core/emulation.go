package usm64core

// MARK: Error

type EmulationError interface{}

// MARK: Context

type EmulationContext struct {
	Registers map[Register]uint64
}

func NewEmulationContext() EmulationContext {
	return EmulationContext{
		Registers: make(map[Register]uint64),
	}
}

type Emulateable interface {
	Emulate(ctx *EmulationContext) EmulationError
}
