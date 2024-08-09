package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type BlockNode[NodeT Node] struct {
	Nodes []NodeT
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

func (p BlockParser[NodeT]) Parse(v *TokenView) (nodes BlockNode[NodeT], err ParsingError) {
	_, err = v.ConsumeToken(lex.LeftCurlyBraceToken)

	if err != nil {
		first, err := p.Parser.Parse(v)
		if err != nil {
			return nodes, GenericUnexpectedError{Expected: p.String()}
		}

		return BlockNode[NodeT]{[]NodeT{first}}, nil
	}

	nodes = BlockNode[NodeT]{ParseMany(p.Parser, v)}
	_, err = v.ConsumeToken(lex.RightCurlyBraceToken)
	return
}
