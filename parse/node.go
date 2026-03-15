package parse

import (
	"slices"
	"strings"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type StringContext struct {
	core.SourceContext
	Indent int
}

// indent returns a string of tabs matching the current indentation level.
func (ctx *StringContext) indent() string {
	return strings.Repeat("\t", ctx.Indent)
}

// renderComments renders a slice of comments as indented whole-line strings.
func (ctx *StringContext) renderComments(comments []lex.Comment) string {
	prefix := ctx.indent()
	var s string
	for _, c := range comments {
		s += prefix + string(c.View.Raw(ctx.SourceContext)) + "\n"
	}
	return s
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
