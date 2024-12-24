package gen

type StepInfo interface{}

// Regular control flow: just continue to execute the instruction that is
// located right after this one in the source code.
type ContinueToNextInstruction struct{}

type ReturnFromFunction struct{}

type JumpToLabel[InstT BaseInstruction] struct {
	Label *LabelInfo[InstT]
}
