package parse

import (
	"usm/lex"
	"usm/source"
)

type TypeNode struct {
	source.UnmanagedSourceView
}

func (n TypeNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n TypeNode) String(ctx source.SourceContext) string {
	return string(n.UnmanagedSourceView.Raw(ctx))
}

type TypeParser struct{}

func (TypeParser) Parse(v *TokenView) (node TypeNode, err ParsingError) {
	tkn, err := v.ConsumeToken(lex.TypToken)
	if err != nil {
		return
	}

	node = TypeNode{tkn.View}
	return
}
