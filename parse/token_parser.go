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

type TokenParser[NodeT TokenNode] struct {
	Token lex.TokenType
}

func (p TokenParser[NodeT]) Parse(v *TokenView) (node NodeT, err core.Result) {
	tkn, err := v.ConsumeToken(p.Token)
	if err != nil {
		return
	}

	return NodeT{tkn.View}, nil
}
