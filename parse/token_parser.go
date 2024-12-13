package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

// MARK: Node

type TokenNode struct {
	core.UnmanagedSourceView
}

func (n TokenNode) View() core.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n TokenNode) String(ctx *StringContext) string {
	return string(n.UnmanagedSourceView.Raw(ctx.SourceContext))
}

// MARK: Parser

type TokenParser[NodeT Node] struct {
	Token       lex.TokenType
	NodeCreator func(lex.Token) NodeT
}

func (p TokenParser[NodeT]) Parse(v *TokenView) (node NodeT, err core.Result) {
	tkn, err := v.ConsumeToken(p.Token)
	if err != nil {
		return
	}

	return p.NodeCreator(tkn), nil
}
