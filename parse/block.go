package parse

import (
	"alon.kr/x/usm/lex"
)

type BlockNode[NodeT any] struct {
	Nodes []NodeT
}

type BlockParser[NodeT any] struct {
	parser Parser[NodeT]
}

func (p BlockParser[NodeT]) String() string {
	return p.parser.String() + " block"
}

func (p BlockParser[NodeT]) Parse(v *TokenView) (node BlockNode[NodeT], err ParsingError) {
	_, err = v.ConsumeToken(lex.LeftCurlyBraceToken)

	if err != nil {
		first, err := p.parser.Parse(v)
		if err != nil {
			return node, GenericUnexpectedError{Expected: p.String()}
		}

		node.Nodes = []NodeT{first}
		return node, nil
	}

	node.Nodes = ParseMany(p.parser, v)
	_, err = v.ConsumeToken(lex.RightCurlyBraceToken)
	return
}
