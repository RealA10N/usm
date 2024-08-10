package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type BlockNode[NodeT Node] struct {
	source.UnmanagedSourceView
	Nodes []NodeT
}

func (n BlockNode[NodeT]) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n BlockNode[NodeT]) String(ctx source.SourceContext) (s string) {
	if len(n.Nodes) == 0 {
		return "{ }\n"
	}

	s = "{\n"
	for _, node := range n.Nodes {
		s += node.String(ctx)
	}
	s += "}\n"

	return
}

type BlockParser[NodeT Node] struct {
	Parser Parser[NodeT]
}

func (p BlockParser[NodeT]) String() string {
	return p.Parser.String() + " block"
}

func (p BlockParser[NodeT]) Parse(v *TokenView) (block BlockNode[NodeT], err ParsingError) {
	leftCurly, err := v.ConsumeToken(lex.LeftCurlyBraceToken)

	if err != nil {
		return
	}

	block.Start = leftCurly.View.Start
	block.Nodes = ParseMany(p.Parser, v)

	rightCurly, err := v.ConsumeTokenIgnoreSeparator(lex.RightCurlyBraceToken)
	if err != nil {
		return
	}

	block.End = rightCurly.View.End
	return block, nil
}
