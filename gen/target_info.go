package gen

import "alon.kr/x/usm/core"

type TargetInfo interface {
	// The location where the target appears in the source code.
	Declaration() *core.UnmanagedSourceView

	// Returns the target string, as it should appear in the source code.
	String() string

	// OnAttach is called when this target is attached to an instruction.
	// Implementations should update any back-references as necessary
	// (e.g. adding the instruction to a register's Definitions list).
	OnAttach(instruction *InstructionInfo)

	// OnDetach is called when this target is detached from an instruction.
	// Implementations should undo the back-references established in OnAttach
	// (e.g. removing the instruction from a register's Definitions list).
	OnDetach(instruction *InstructionInfo)
}
