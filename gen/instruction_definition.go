package gen

import (
	"alon.kr/x/usm/core"
)

type BaseInstruction interface {
	// This method is usd by the USM engine to generate the internal control
	// flow graph representation.
	//
	// Should return a non-empty slice. If the instruction does not have any
	// consecutive steps in the function (for example, a return statement),
	// then a special dedicated return step should be returned.
	PossibleNextSteps() (StepInfo, core.ResultList)

	// Returns the string that represents the operator of the instruction.
	// For example, for the add instruction this method would return "ADD".
	String() string
}

// A basic instruction definition. This defines the logic that converts the
// generic, architecture / instruction set independent instruction AST nodes
// into a format instruction which is part of a specific instruction set.
type InstructionDefinition interface {
	// Build an instruction from the provided instruction information.
	BuildInstruction(info *InstructionInfo) (BaseInstruction, core.ResultList)
}
