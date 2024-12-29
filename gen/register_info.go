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

	// The first location in the source code in which the register is declared
	// or assigned a value.
	Declaration core.UnmanagedSourceView
}

func (i *RegisterInfo) AddDefinition(info *InstructionInfo) {
	i.Definitions = append(i.Definitions, info)
}

func (i *RegisterInfo) toPartialRegisterInfo() registerPartialInfo {
	return registerPartialInfo{
		Name:        i.Name,
		Type:        &i.Type,
		Declaration: i.Declaration,
	}
}
