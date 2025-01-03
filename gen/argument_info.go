package gen

import "alon.kr/x/usm/core"

type ArgumentInfo interface {
	// A pointer to the ReferencedTypeInfo instance that corresponds to the
	// type of the argument. Nil if the argument does not have a type (for
	// example, a label).
	GetType() *ReferencedTypeInfo

	// The location where the argument appears in the source code.
	Declaration() core.UnmanagedSourceView

	// Returns the argument string, as it should appear in the code code.
	String() string
}
