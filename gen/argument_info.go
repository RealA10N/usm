package gen

import "alon.kr/x/usm/core"

type ArgumentInfo interface {
	// The location where the argument appears in the source code.
	Declaration() *core.UnmanagedSourceView

	// Returns the argument string, as it should appear in the code code.
	String() string
}
