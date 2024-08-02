package parse

import (
	"usm/lex"
)

type ArgumentNode struct {
	Node
	Type     lex.Token
	Register lex.Token
}

type ArgumentNodeParser struct{}

func (p ArgumentNodeParser) Parse(v *TokenView) (node ArgumentNode, err ParsingError) {
	typ, err := ConsumeToken(v, lex.TypToken)
	if err != nil {
		return
	}

	reg, err := ConsumeToken(v, lex.RegToken)
	if err != nil {
		return
	}

	node = ArgumentNode{
		Node:     NewNodeFromBoundaryTokens(typ, reg),
		Type:     typ,
		Register: reg,
	}

	return node, nil
}
