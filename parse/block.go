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
		return "{ }"
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
