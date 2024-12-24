package gen

import "alon.kr/x/usm/core"

type InstructionManager[InstT BaseInstruction] interface {
	// Get the instruction definition that corresponds to the provided name.
	GetInstructionDefinition(name string) (InstructionDefinition[InstT], core.ResultList)
}
