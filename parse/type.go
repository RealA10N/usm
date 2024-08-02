package parse

import (
	"usm/lex"
	"usm/source"
)

type TypeNode struct {
	view source.UnmanagedSourceView
}

func (n TypeNode) View() source.UnmanagedSourceView {
	return n.view
}

func (n TypeNode) String(ctx source.SourceContext) string {
	return string(n.view.Raw(ctx))
}

type TypeParser struct{}

func (TypeParser) Parse(v *TokenView) (node TypeNode, err ParsingError) {
	tkn, err := ConsumeToken(v, lex.TypToken)
	if err != nil {
		return
	}

	node = TypeNode{view: tkn.View}
	return
}
