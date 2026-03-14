package parse

import (
	"slices"
	"sort"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type StringContext struct {
	core.SourceContext
	Indent   int
	Comments []lex.Comment // sorted by View.Start; never mutated
	cursor   int           // index of next comment to process
}

// WholeLineCommentsBefore returns all unprocessed comments that appear before
// nodeStart in the source. Advances the internal cursor.
func (ctx *StringContext) WholeLineCommentsBefore(nodeStart core.SourceViewOffset) []lex.Comment {
	hi := sort.Search(len(ctx.Comments), func(i int) bool {
		return ctx.Comments[i].View.Start >= nodeStart
	})
	result := ctx.Comments[ctx.cursor:hi]
	ctx.cursor = hi
	return result
}

// InlineComment returns the trailing comment on the same source line as nodeEnd,
// if one exists. Advances the cursor past it.
func (ctx *StringContext) InlineComment(nodeEnd core.SourceViewOffset) string {
	if ctx.cursor >= len(ctx.Comments) {
		return ""
	}
	lineEnd := core.SourceViewOffset(len(ctx.SourceContext))
	for i := int(nodeEnd); i < len(ctx.SourceContext); i++ {
		if ctx.SourceContext[i] == '\n' {
			lineEnd = core.SourceViewOffset(i)
			break
		}
	}
	c := ctx.Comments[ctx.cursor]
	if c.View.Start >= nodeEnd && c.View.Start < lineEnd {
		ctx.cursor++
		return " " + string(c.View.Raw(ctx.SourceContext))
	}
	return ""
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
