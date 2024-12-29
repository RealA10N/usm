package gen

import "alon.kr/x/usm/core"

type InstructionInfo struct {
	// The targets of the instruction.
	Targets []*RegisterArgumentInfo

	// The arguments of the instruction.
	Arguments []ArgumentInfo

	// The labels that point directly to this instruction.
	Labels []*LabelInfo

	// The actual instruction information, which is ISA specific.
	Instruction BaseInstruction

	// The location in which the instruction was defined in the source code.
	// Can be nil if the instruction was defined internally, for example,
	// in an optimization.
	Declaration *core.UnmanagedSourceView
}

func NewEmptyInstructionInfo(
	declaration *core.UnmanagedSourceView,
) *InstructionInfo {
	return &InstructionInfo{
		Targets:     []*RegisterArgumentInfo{},
		Arguments:   []ArgumentInfo{},
		Labels:      []*LabelInfo{},
		Instruction: nil,
		Declaration: declaration,
	}
}
