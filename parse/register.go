package parse

import "alon.kr/x/usm/lex"

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
