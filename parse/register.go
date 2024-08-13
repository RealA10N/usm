package parse

import "alon.kr/x/usm/lex"

type RegisterNode = TokenNode

func NewRegisterParser() Parser[RegisterNode] {
	return TokenParser[RegisterNode]{lex.RegisterToken}
}
