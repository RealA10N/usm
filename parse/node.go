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
	Comments []lex.Comment // sorted by View.Start; never mutated; nil = no comments
}

// CommentsInRange returns all comments with Start in [start, end).
// Uses binary search on the sorted Comments slice.
// Used for whole-line comments where multiple comments may appear in a gap.
func (ctx *StringContext) CommentsInRange(start, end core.SourceViewOffset) []lex.Comment {
	lo := sort.Search(len(ctx.Comments), func(i int) bool {
		return ctx.Comments[i].View.Start >= start
	})
	hi := sort.Search(len(ctx.Comments), func(i int) bool {
		return ctx.Comments[i].View.Start >= end
	})
	return ctx.Comments[lo:hi]
}

// WholeLineCommentsAfter returns whole-line comments that appear between
// prevNodeEnd and nextNodeStart. It skips past the first '\n' after prevNodeEnd
// so that inline comments on the same line as the previous node are excluded
// (those are handled by each node's own String() method via InlineComment).
func (ctx *StringContext) WholeLineCommentsAfter(prevNodeEnd, nextNodeStart core.SourceViewOffset) []lex.Comment {
	start := prevNodeEnd
	for i := int(prevNodeEnd); i < len(ctx.SourceContext); i++ {
		if ctx.SourceContext[i] == '\n' {
			start = core.SourceViewOffset(i + 1)
			break
		}
	}
	return ctx.CommentsInRange(start, nextNodeStart)
}

// InlineComment returns the comment (if any) that starts after nodeEnd but
// before the next '\n' in the source — i.e., on the same original source line.
// Returns "" if none. At most one comment can exist per source line.
func (ctx *StringContext) InlineComment(nodeEnd core.SourceViewOffset) string {
	lineEnd := core.SourceViewOffset(len(ctx.SourceContext))
	for i := int(nodeEnd); i < len(ctx.SourceContext); i++ {
		if ctx.SourceContext[i] == '\n' {
			lineEnd = core.SourceViewOffset(i)
			break
		}
	}
	lo := sort.Search(len(ctx.Comments), func(i int) bool {
		return ctx.Comments[i].View.Start >= nodeEnd
	})
	if lo < len(ctx.Comments) && ctx.Comments[lo].View.Start < lineEnd {
		return " " + string(ctx.Comments[lo].View.Raw(ctx.SourceContext))
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
