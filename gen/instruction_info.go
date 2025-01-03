package gen

import "alon.kr/x/usm/core"

type InstructionInfo struct {
	*BasicBlockInfo

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
		BasicBlockInfo: nil,
		Targets:        []*RegisterArgumentInfo{},
		Arguments:      []ArgumentInfo{},
		Labels:         []*LabelInfo{},
		Instruction:    nil,
		Declaration:    declaration,
	}
}

func (i *InstructionInfo) LinkToBasicBlock(basicBlock *BasicBlockInfo) {
	i.BasicBlockInfo = basicBlock
	for _, label := range i.Labels {
		label.LinkToBasicBlock(basicBlock)
	}
}

func (i *InstructionInfo) String() string {
	s := ""
	for _, label := range i.Labels {
		s += label.String() + "\n"
	}

	s += "\t"

	if len(i.Targets) > 0 {
		for _, target := range i.Targets {
			s += target.String() + " "
		}
		s += "= "
	}

	s += i.Instruction.String()

	if len(i.Arguments) > 0 {
		for _, argument := range i.Arguments {
			s += " " + argument.String()
		}
	}

	return s
}
