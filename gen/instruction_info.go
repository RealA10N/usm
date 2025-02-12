package gen

import "alon.kr/x/usm/core"

type InstructionInfo struct {
	*BasicBlockInfo

	// The targets of the instruction.
	// TODO: is a pointer reference really required here?
	Targets []*TargetInfo

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
		Targets:        []*TargetInfo{},
		Arguments:      []ArgumentInfo{},
		Labels:         []*LabelInfo{},
		Instruction:    nil,
		Declaration:    declaration,
	}
}

func (i *InstructionInfo) linkToBasicBlock(basicBlock *BasicBlockInfo) {
	i.BasicBlockInfo = basicBlock
	for _, label := range i.Labels {
		label.linkToBasicBlock(basicBlock)
	}
}

// Appends the given register(s) as a target(s) of the instruction,
// including updating the required instruction and register information fields.
func (i *InstructionInfo) AppendTarget(targets ...*TargetInfo) {
	for _, target := range targets {
		target.Register.AddDefinition(i)
	}

	i.Targets = append(i.Targets, targets...)
}

func (i *InstructionInfo) AppendArgument(arguments ...ArgumentInfo) {
	i.Arguments = append(i.Arguments, arguments...)
}

func (i *InstructionInfo) AppendLabels(labels ...*LabelInfo) {
	i.Labels = append(i.Labels, labels...)
	for _, label := range labels {
		label.linkToBasicBlock(i.BasicBlockInfo)
	}
}

func (i *InstructionInfo) MoveLabels(targetInstruction *InstructionInfo) {
	targetInstruction.AppendLabels(i.Labels...)
	i.Labels = nil
}

// Updates the internal instruction instance to the provided one.
//
// This can be used to update the instruction, but keep the same arguments and
// targets, for example, as an optimization to a more specific operation which
// accepts the same arguments in certain cases.
func (i *InstructionInfo) SetBaseInstruction(instruction BaseInstruction) {
	i.Instruction = instruction
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
