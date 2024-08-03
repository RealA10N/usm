package parse

import "usm/lex"

type Parser[Node any] interface {
	Parse(v *TokenView) (Node, ParsingError)
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

func ParseManyConsumeSeperators[Node any](p Parser[Node], v *TokenView) (nodes []Node) {
	for {
		v.ConsumeManyTokens(lex.SepToken)
		inst, err := p.Parse(v)
		if err != nil {
			break
		}
		nodes = append(nodes, inst)
	}
	v.ConsumeManyTokens(lex.SepToken)
	return nodes
}
