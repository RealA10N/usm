package parse

import (
	"usm/source"
)

type Node interface {
	// Return a reference to the node substring in the source code
	View() source.UnmanagedSourceView

	// Regenerate ("format") the code to a unique, single representation.
	String(ctx source.SourceContext) string
}
