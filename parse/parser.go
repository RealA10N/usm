package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type Parser[Node any] interface {
	Parse(v *TokenView) (Node, core.Result)
}

func ParseMany[Node any](p Parser[Node], v *TokenView) (nodes []Node) {
	for {
		typ, err := p.Parse(v)
		if err != nil {
			return
		}
		nodes = append(nodes, typ)
	}
}

func ParseManyIgnoreSeparators[Node any](
	p Parser[Node],
	v *TokenView,
) (nodes []Node, err core.Result) {
	v.ConsumeManyTokens(lex.SeparatorToken)
	for {
		var node Node
		node, err = p.Parse(v)
		if err != nil {
			return
		}
		nodes = append(nodes, node)

		v.ConsumeManyTokens(lex.SeparatorToken)
	}
}
