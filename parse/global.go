package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type GlobalNode struct {
	source.UnmanagedSourceView
}

func (n GlobalNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n GlobalNode) String(ctx source.SourceContext) string {
	return string(n.UnmanagedSourceView.Raw(ctx.ViewContext))
}

type GlobalParser struct{}

func (GlobalParser) Parse(v *TokenView) (node GlobalNode, err ParsingError) {
	tkn, err := v.ConsumeToken(lex.GlobalToken)
	if err != nil {
		return
	}

	return GlobalNode{tkn.View}, nil
}
