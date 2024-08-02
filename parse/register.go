package parse

import (
	"usm/lex"
	"usm/source"
)

type RegisterNode struct {
	view source.UnmanagedSourceView
}

func (n RegisterNode) View() source.UnmanagedSourceView {
	return n.view
}

type RegisterParser struct{}

func (RegisterParser) Parse(v *TokenView) (node RegisterNode, err ParsingError) {
	tkn, err := ConsumeToken(v, lex.RegToken)
	if err != nil {
		return
	}

	node = RegisterNode{view: tkn.View}
	return
}
