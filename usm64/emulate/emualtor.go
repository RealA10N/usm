package usm64emulate

import (
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type Emulator struct{}

func (Emulator) Emulate(
	function *gen.FunctionInfo[usm64core.Instruction],
) usm64core.EmulationError {
	ctx := usm64core.NewEmulationContext()

	for ctx.NextInstructionIndex < uint64(len(function.Instructions)) {
		nextInstruction := function.Instructions[ctx.NextInstructionIndex]
		err := nextInstruction.Emulate(&ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewEmulator() Emulator {
	return Emulator{}
}
