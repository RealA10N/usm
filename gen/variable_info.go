package gen

import "alon.kr/x/usm/core"

// VariableInfo represents a named, typed slot in the function's local frame.
// Variables are distinct from registers: they are mutable memory locations
// (typically stack-allocated by the backend) that can be read and written by
// multiple instructions, unlike registers which are defined exactly once.
type VariableInfo struct {
	// The name of the variable as it appears in the source code, including
	// the '&' sigil (e.g. "&local").
	Name string

	// The declared type of the variable.
	Type ReferencedTypeInfo

	// The location in the source code where the variable was declared.
	Declaration core.UnmanagedSourceView
}

func (i *VariableInfo) String() string {
	return i.Name
}
