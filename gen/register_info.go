package gen

import "alon.kr/x/usm/core"

type RegisterInfo struct {
	// The name of the register, as it appears in the source code.
	Name string

	// The type of the register.
	Type ReferencedTypeInfo

	// Instructions in which the register is a target, and is defined or
	// assigned a new value.
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
	Declaration core.UnmanagedSourceView
}

func (i *RegisterInfo) AddDefinition(info *InstructionInfo) {
	i.Definitions = append(i.Definitions, info)
}

func (i *RegisterInfo) AddUsage(info *InstructionInfo) {
	i.Usages = append(i.Usages, info)
}

func (i *RegisterInfo) toPartialRegisterInfo() registerPartialInfo {
	return registerPartialInfo{
		Name:        i.Name,
		Type:        &i.Type,
		Declaration: i.Declaration,
	}
}
