package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type InstructionManager[InstT BaseInstruction] interface {
	// Get the instruction definition that corresponds to the provided name.
	//
	// Instruction node is for extra context, if needed, especially for
	// generating nice error messages.
	GetInstructionDefinition(
		name string,
		node parse.InstructionNode,
	) (InstructionDefinition[InstT], core.ResultList)
}
