package gen

import "alon.kr/x/usm/core"

type ArgumentInfo interface {
	// The location where the argument appears in the source code.
	Declaration() *core.UnmanagedSourceView

	// Returns the argument string, as it should appear in the code code.
	String() string

	// OnAttach is called when this argument is attached to an instruction.
	// Implementations should update any back-references as necessary
	// (e.g. adding the instruction to a register's Usages list).
	OnAttach(instruction *InstructionInfo)

	// OnDetach is called when this argument is detached from an instruction.
	// Implementations should undo the back-references established in OnAttach
	// (e.g. removing the instruction from a register's Usages list).
	OnDetach(instruction *InstructionInfo)
}
