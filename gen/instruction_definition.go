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
	// Validation runs as a read-only pass after IR generation is complete.
	// Implementations must not mutate any IR nodes or infer missing information —
	// those concerns belong in the generation phase. Validate should only inspect
	// the already-built IR and return errors if it is inconsistent.
	Validate(*InstructionInfo) core.ResultList
}

type NonBranchingInstruction struct{}

func (NonBranchingInstruction) PossibleNextSteps(*InstructionInfo) (StepInfo, core.ResultList) {
	return StepInfo{
		PossibleContinue: true,
	}, core.ResultList{}
}

type BranchToLabelArguments struct{}

func (BranchToLabelArguments) PossibleNextSteps(info *InstructionInfo) (StepInfo, core.ResultList) {
	return StepInfo{
		PossibleBranches: ArgumentsToLabels(info.Arguments),
	}, core.ResultList{}
}

type BranchesToLabelArgumentsOrContinues struct{}

func (BranchesToLabelArgumentsOrContinues) PossibleNextSteps(info *InstructionInfo) (StepInfo, core.ResultList) {
	return StepInfo{
		PossibleBranches: ArgumentsToLabels(info.Arguments),
		PossibleContinue: true,
	}, core.ResultList{}
}

type ReturningInstruction struct{}

func (ReturningInstruction) PossibleNextSteps(*InstructionInfo) (StepInfo, core.ResultList) {
	return StepInfo{
		PossibleReturn: true,
	}, core.ResultList{}
}
