package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type TokenNode struct {
	source.UnmanagedSourceView
}

func (n TokenNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n TokenNode) String(ctx *StringContext) string {
	return string(n.UnmanagedSourceView.Raw(ctx.SourceContext))
}

type TokenParser[NodeT TokenNode] struct {
	Token lex.TokenType
}

func (p TokenParser[NodeT]) Parse(v *TokenView) (node NodeT, err ParsingError) {
	tkn, err := v.ConsumeToken(p.Token)
	if err != nil {
		return
	}

	return NodeT{tkn.View}, nil
}
