package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type RegisterNode struct{ TokenNode }
type RegisterParser struct{ TokenParser[RegisterNode] }

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

// AnnotatedRegisterParser parses a register token, optionally preceded by
// a type annotation which is accepted but discarded (the register already
// carries its declared type).
type AnnotatedRegisterParser struct {
	TypeParser
	RegisterParser Parser[RegisterNode]
}

func NewAnnotatedRegisterParser() Parser[RegisterNode] {
	return AnnotatedRegisterParser{
		RegisterParser: NewRegisterParser(),
	}
}

func (p AnnotatedRegisterParser) Parse(v *TokenView) (RegisterNode, core.Result) {
	// If the next two tokens are ($type, %register), consume the type annotation.
	if _, err := v.PeekToken(lex.TypeToken); err == nil {
		if tkn1, err := v.At(1); err == nil && tkn1.Type == lex.RegisterToken {
			p.TypeParser.Parse(v) // consume and discard
		}
	}
	return p.RegisterParser.Parse(v)
}
