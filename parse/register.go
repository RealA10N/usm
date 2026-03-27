package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type RegisterNode struct {
	TokenNode

	// Optional explicit type annotation (e.g. "$32" in "$32 %a").
	// Set whenever a type token immediately precedes the register token,
	// whether the register appears as an instruction target or as an argument.
	Type *TypeNode
}

func (n RegisterNode) View() core.UnmanagedSourceView {
	v := n.TokenNode.View()
	if n.Type != nil {
		v = v.MergeStart(n.Type.View())
	}
	return v
}

type RegisterParser struct {
	TypeParser
	TokenParser[RegisterNode]
}

func RegisterNodeCreator(tkn lex.Token) RegisterNode {
	return RegisterNode{TokenNode: TokenNode{tkn.View}}
}

func (n RegisterNode) String(ctx *StringContext) string {
	if n.Type != nil {
		return n.Type.String(ctx) + " " + n.TokenNode.String(ctx)
	}
	return n.TokenNode.String(ctx)
}

func NewRegisterParser() Parser[RegisterNode] {
	return RegisterParser{
		TokenParser: TokenParser[RegisterNode]{
			Token:       lex.RegisterToken,
			NodeCreator: RegisterNodeCreator,
		},
	}
}

// Parse optionally consumes a preceding type annotation, then parses the
// register token. This allows both "%name" and "$type %name" syntax in
// argument position.  The type annotation is preserved in RegisterNode.Type.
func (p RegisterParser) Parse(v *TokenView) (RegisterNode, core.Result) {
	var typ *TypeNode
	if v.PeekTokens(lex.TypeToken, lex.RegisterToken) {
		t, err := p.TypeParser.Parse(v)
		if err != nil {
			return RegisterNode{}, err
		}
		typ = &t
	}
	node, err := p.TokenParser.Parse(v)
	node.Type = typ
	return node, err
}
