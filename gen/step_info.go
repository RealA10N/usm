package gen

type StepInfo struct {
	// A list of all branches that the execution might jump to after the
	// execution of the instruction.
	PossibleBranches []*LabelInfo

	// True if it is possible that execution after this instruction should
	// continue to the next instruction in the function.
	PossibleContinue bool

	// True if it is possible that execution after this instruction should
	// terminate the function execution. This can be used for termination of
	// the entire program, or a return from a function.
	PossibleReturn bool
}

func (p *StepInfo) IsBranchPossible() bool {
	return len(p.PossibleBranches) > 0
}

func (p *StepInfo) DefinitelyContinue() bool {
	return p.PossibleContinue && !p.IsBranchPossible() && !p.PossibleReturn
}

func (p *StepInfo) DefinitelyReturn() bool {
	return p.PossibleReturn && !p.IsBranchPossible() && !p.PossibleContinue
}
