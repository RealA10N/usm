package gen

import "alon.kr/x/usm/core"

type GlobalInfo interface {
	Name() string
	Declaration() *core.UnmanagedSourceView

	// IsDefined returns true if the global variable is defined in the current
	// file. If false, it means that the global is externally defined and
	// expected to be linked in a later stage.
	IsDefined() bool
}
