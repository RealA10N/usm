package parse

import (
	"strings"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type BlockNode[NodeT Node] struct {
	core.UnmanagedSourceView
	Nodes []NodeT
}

func (n BlockNode[NodeT]) View() core.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n BlockNode[NodeT]) String(ctx *StringContext) (s string) {
	if len(n.Nodes) == 0 {
		// Consume any inline comment after '{' and any whole-line comments inside
		// the block so they don't leak to the outer scope.
		inline := ctx.InlineComment(n.UnmanagedSourceView.Start + 1)
		inner := ctx.WholeLineCommentsBefore(n.UnmanagedSourceView.End)
		if inline == nil && len(inner) == 0 {
			return "{ }"
		}
		// The block has comments but no instructions; expand to multi-line.
		prefix := strings.Repeat("\t", ctx.Indent+1)
		s = "{\n"
		if inline != nil {
			s += prefix + string(inline.View.Raw(ctx.SourceContext)) + "\n"
		}
		for _, c := range inner {
			s += prefix + string(c.View.Raw(ctx.SourceContext)) + "\n"
		}
		return s + strings.Repeat("\t", ctx.Indent) + "}"
	}

	s = "{\n"
	ctx.Indent++
	for _, node := range n.Nodes {
		s += node.String(ctx)
	}
	// Emit any comments between the last node and the closing '}'.
	prefix := strings.Repeat("\t", ctx.Indent)
	for _, c := range ctx.WholeLineCommentsBefore(n.UnmanagedSourceView.End) {
		s += prefix + string(c.View.Raw(ctx.SourceContext)) + "\n"
	}
	ctx.Indent--
	s += strings.Repeat("\t", ctx.Indent) + "}"

	return
}

type BlockParser[NodeT Node] struct {
	Parser Parser[NodeT]
}

func (p BlockParser[NodeT]) Parse(v *TokenView) (block BlockNode[NodeT], err core.Result) {
	leftCurly, err := v.ConsumeToken(lex.LeftCurlyBraceToken)

	if err != nil {
		return
	}

	block.Start = leftCurly.View.Start
	block.Nodes, _ = ParseManyIgnoreSeparators(p.Parser, v)

	rightCurly, err := v.ConsumeTokenIgnoreSeparator(lex.RightCurlyBraceToken)
	if err != nil {
		return
	}

	block.End = rightCurly.View.End
	return block, nil
}
