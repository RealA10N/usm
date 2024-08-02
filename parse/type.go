package parse

import (
	"usm/lex"
	"usm/source"
)

type TypeNode struct {
	View source.UnmanagedSourceView
}

type TypeParser struct{}

func (p TypeParser) Parse(v *TokenView) (node TypeNode, err ParsingError) {
	tkn, err := ConsumeToken(v, lex.TypToken)
	if err != nil {
		return
	}

	node = TypeNode{View: tkn.View}
	return
}
