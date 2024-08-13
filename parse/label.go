package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type LabelNode struct {
	source.UnmanagedSourceView
}

func (n LabelNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n LabelNode) String(ctx *StringContext) string {
	return string(n.UnmanagedSourceView.Raw(ctx.SourceContext))
}

type LabelParser struct{}

func (LabelParser) Parse(v *TokenView) (node LabelNode, err ParsingError) {
	tkn, err := v.ConsumeToken(lex.LabelToken)
	if err != nil {
		return
	}

	return LabelNode{tkn.View}, nil
}
