package parse

import (
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

// stringEmptyBlock formats a block that has no instruction nodes.
// Any comments inside the block range are preserved; if present, the block
// expands to multi-line to house them. Comments after '}' are left unconsumed
// for the outer context (e.g. an inline comment on the same line as '}').
func (n BlockNode[NodeT]) stringEmptyBlock(ctx *StringContext) string {
	inner := ctx.FormatCommentsBeforeIndented(n.UnmanagedSourceView.End, ctx.indent()+"\t")
	if inner == "" {
		return "{ }"
	}
	return "{\n" + inner + ctx.indent() + "}"
}

func (n BlockNode[NodeT]) String(ctx *StringContext) (s string) {
	if len(n.Nodes) == 0 {
		return n.stringEmptyBlock(ctx)
	}

	s = "{\n"
	ctx.Indent++
	for _, node := range n.Nodes {
		s += node.String(ctx)
	}
	s += ctx.FormatCommentsBeforeIndented(n.UnmanagedSourceView.End, ctx.indent())
	ctx.Indent--
	s += ctx.indent() + "}"
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
