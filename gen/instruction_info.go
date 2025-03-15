package gen

import "alon.kr/x/usm/core"

type InstructionInfo struct {
	*BasicBlockInfo

	// The targets of the instruction.
	// TODO: is a pointer reference really required here?
	Targets []*TargetInfo

	// The arguments of the instruction.
	Arguments []ArgumentInfo

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
		Instruction:    nil,
		Declaration:    declaration,
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

func (i *InstructionInfo) SwitchTarget(
	target *TargetInfo,
	newRegister *RegisterInfo,
) {
	target.Register.RemoveDefinition(i)
	target.Register = newRegister
	target.Register.AddDefinition(i)
}

func (i *InstructionInfo) AppendArgument(arguments ...ArgumentInfo) {
	i.Arguments = append(i.Arguments, arguments...)
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
	operator := i.Instruction.String()

	if len(i.Targets) > 0 {
		for _, target := range i.Targets {
			s += target.String() + " "
		}
		s += "="

		if len(operator) > 0 {
			s += " "
		}
	}

	s += operator

	if len(i.Arguments) > 0 {
		for _, argument := range i.Arguments {
			s += " " + argument.String()
		}
	}

	return s
}
