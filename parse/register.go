package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type RegisterNode struct {
	source.UnmanagedSourceView
}

func (n RegisterNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n RegisterNode) String(ctx source.SourceContext) string {
	return string(n.UnmanagedSourceView.Raw(ctx))
}

type RegisterParser struct{}

func (RegisterParser) Parse(v *TokenView) (node RegisterNode, err ParsingError) {
	tkn, err := v.ConsumeToken(lex.RegToken)
	if err != nil {
		return
	}

	node = RegisterNode{tkn.View}
	return
}
