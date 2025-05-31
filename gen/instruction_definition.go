package gen

import (
	"alon.kr/x/usm/core"
)

type InstructionDefinition interface {
	// This method is usd by the USM engine to generate the internal control
	// flow graph representation.
	//
	// Should return a non-empty slice. If the instruction does not have any
	// consecutive steps in the function (for example, a return statement),
	// then a special dedicated return step should be returned.
	PossibleNextSteps(*InstructionInfo) (StepInfo, core.ResultList)

	// Returns the string that represents the operator of the instruction.
	// For example, for the add instruction this method would return "ADD".
	//
	// This is required because some instructions may be generated automatically,
	// and we want to be able to display them in a human-readable format.
	Operator(*InstructionInfo) string

	// Validate the instruction information structure, according to the
	// expected arguments, targets, and other related information.
	//
	// Instruction validation should as one of the last steps in the compilation
	// process, and you should be able to assume that all relevant information
	// in the structures is filled in and propagated correctly.
	Validate(*InstructionInfo) core.ResultList
}

type NonBranchingInstruction struct{}

func (NonBranchingInstruction) PossibleNextSteps(*InstructionInfo) (StepInfo, core.ResultList) {
	return StepInfo{
		PossibleContinue: true,
	}, core.ResultList{}
}

type ReturningInstruction struct{}

func (ReturningInstruction) PossibleNextSteps(*InstructionInfo) (StepInfo, core.ResultList) {
	return StepInfo{
		PossibleReturn: true,
	}, core.ResultList{}
}
