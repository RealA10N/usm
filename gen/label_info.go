package gen

import "alon.kr/x/usm/core"

type LabelInfo struct {
	// The name of the label, as it appears in the source code.
	Name string

	// The index of the instruction that the label is referencing.
	InstructionIndex core.UsmUint

	// A view of the label declaration in the source code.
	Declaration core.UnmanagedSourceView
}
