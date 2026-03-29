package gen

import (
	"slices"

	"alon.kr/x/usm/core"
)

type RegisterInfo struct {
	// The name of the register, as it appears in the source code.
	Name string

	// The type of the register.
	Type ReferencedTypeInfo

	// All instructions that reference this register, either as a target or as
	// an argument. An instruction appears once per argument/target position that
	// references this register (e.g. "add %x %x" adds two entries).
	//
	// Note: this list does not include references via function parameters.
	//
	// TODO: for quicker updates of the data structure, this field should be a
	// linked list where each entry points to the exact target/argument slot in
	// the relevant instruction, not to the whole instruction struct.
	References []*InstructionInfo

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

func (i *RegisterInfo) AddReference(info *InstructionInfo) {
	i.References = append(i.References, info)
}

// RemoveReference removes one occurrence of info from the References list.
// If the same instruction appears multiple times (e.g. "add %x %x"), each
// call removes exactly one entry, matching one argument/target position being
// detached.
func (i *RegisterInfo) RemoveReference(info *InstructionInfo) {
	if idx := slices.Index(i.References, info); idx != -1 {
		i.References = slices.Delete(i.References, idx, idx+1)
	}
}
