package gen

import "alon.kr/x/usm/core"

type RegisterInfo struct {
	// The name of the register, as it appears in the source code.
	Name string

	// The type of the register.
	Type ReferencedTypeInfo

	// Instructions in which the register is a target, and is defined or
	// assigned a new value.
	//
	// Note: This list is not a complete representation of all locations in which
	// the register is defined, since it can be defined as a function parameter.
	Definitions []*InstructionInfo

	// Instructions in which the register appears as a source, i.e. as an
	// read only argument.
	Usages []*InstructionInfo

	// TODO: for quicker updates of the data structure, both `Definitions` and
	// `Usages` fields should be a linked list, where each entry points to the
	// exact target/argument in the relevant instruction, and not to the whole
	// instruction information struct.

	// The first location in the source code in which the register is declared
	// or assigned a value.
	// TODO: make this not required, for generated registers.
	Declaration core.UnmanagedSourceView
}

func NewRegisterInfo(name string, typ ReferencedTypeInfo) *RegisterInfo {
	return &RegisterInfo{
		Name: name,
		Type: typ,
	}
}

// Return the string that represents the register, as it should appear in the
// source code.
func (i *RegisterInfo) String() string {
	return i.Name
}

func (i *RegisterInfo) AddDefinition(info *InstructionInfo) {
	i.Definitions = append(i.Definitions, info)
}

func (i *RegisterInfo) AddUsage(info *InstructionInfo) {
	i.Usages = append(i.Usages, info)
}
