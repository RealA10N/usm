package gen

import "alon.kr/x/usm/core"

type InstructionInfo struct {
	*BasicBlockInfo

	// The targets of the instruction.
	Targets []TargetInfo

	// The arguments of the instruction.
	Arguments []ArgumentInfo

	// The actual instruction information, which is ISA specific.
	Definition InstructionDefinition

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
		Targets:        []TargetInfo{},
		Arguments:      []ArgumentInfo{},
		Definition:     nil,
		Declaration:    declaration,
	}
}

func (i *InstructionInfo) Validate() core.ResultList {
	return i.Definition.Validate(i)
}

// Appends the given target(s) to the instruction,
// including updating the required instruction and register information fields.
func (i *InstructionInfo) AppendTarget(targets ...TargetInfo) {
	for _, target := range targets {
		target.OnAttach(i)
	}

	i.Targets = append(i.Targets, targets...)
}

func (i *InstructionInfo) AppendArgument(arguments ...ArgumentInfo) {
	i.Arguments = append(i.Arguments, arguments...)
}

// SubstituteArgument replaces the argument at the given index with newArg,
// keeping any back-references (e.g. register usage lists) consistent by
// delegating to OnDetach on the old argument and OnAttach on the new one.
func (i *InstructionInfo) SubstituteArgument(index int, newArg ArgumentInfo) {
	i.Arguments[index].OnDetach(i)
	newArg.OnAttach(i)
	i.Arguments[index] = newArg
}

// Updates the internal instruction instance to the provided one.
//
// This can be used to update the instruction, but keep the same arguments and
// targets, for example, as an optimization to a more specific operation which
// accepts the same arguments in certain cases.
func (i *InstructionInfo) SetInstruction(instruction InstructionDefinition) {
	i.Definition = instruction
}

func (i *InstructionInfo) String() string {
	s := ""
	operator := i.Definition.Operator(i)

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
