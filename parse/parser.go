package parse

import (
	"fmt"

	"alon.kr/x/usm/lex"
)

type Parser[Node any] interface {
	fmt.Stringer
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

// Parse nodes using the provided parser, but after each node consume at least
// one separator.
//
// This is useful for nodes that expect that they end in a separator (line break),
// such as the InstructionNode or the FunctionNode.
func ParseManyConsumeSeparators[Node any](
	p Parser[Node],
	v *TokenView,
) (nodes []Node, err ParsingError) {
	for {
		var node Node
		node, err = p.Parse(v)
		if err != nil {
			return
		}
		nodes = append(nodes, node)

		err = v.ConsumeAtLeastTokens(1, lex.SeparatorToken)
		if err != nil {
			return
		}
	}
}

func ParseManyIgnoreSeparators[Node any](
	p Parser[Node],
	v *TokenView,
) (nodes []Node, err ParsingError) {
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
