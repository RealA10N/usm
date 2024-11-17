package parse

import (
	"slices"

	"alon.kr/x/usm/core"
)

type StringContext struct {
	core.SourceContext
	Indent int
}

type Node interface {
	// Return a reference to the node substring in the source code
	View() core.UnmanagedSourceView

	// Regenerate ("format") the code to a unique, single representation.
	String(ctx *StringContext) string
}

// This function sorts the nodes according to their source order.
func SortNodesBySourceOrder(nodes []Node) {
	slices.SortFunc(nodes, func(i, j Node) int {
		return int(i.View().Start) - int(j.View().Start)
	})
}
