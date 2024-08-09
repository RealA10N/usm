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

func (RegisterParser) String() string {
	return "%register"
}

func (RegisterParser) Parse(v *TokenView) (node RegisterNode, err ParsingError) {
	tkn, err := v.ConsumeToken(lex.RegisterToken)
	if err != nil {
		return
	}

	return RegisterNode{tkn.View}, nil
}
