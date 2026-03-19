package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type RegisterNode struct{ TokenNode }

type RegisterParser struct {
	TypeParser
	TokenParser[RegisterNode]
}

func RegisterNodeCreator(tkn lex.Token) RegisterNode {
	return RegisterNode{TokenNode{tkn.View}}
}

func NewRegisterParser() Parser[RegisterNode] {
	return RegisterParser{
		TokenParser: TokenParser[RegisterNode]{
			Token:       lex.RegisterToken,
			NodeCreator: RegisterNodeCreator,
		},
	}
}

// Parse optionally consumes a preceding type annotation (discarded), then
// parses the register token. This allows both "%name" and "$type %name"
// syntax in argument position.
func (p RegisterParser) Parse(v *TokenView) (RegisterNode, core.Result) {
	if v.PeekTokens(lex.TypeToken, lex.RegisterToken) {
		p.TypeParser.Parse(v) // consume and discard the type prefix
	}
	return p.TokenParser.Parse(v)
}
