package parse

import "alon.kr/x/usm/lex"

type RegisterNode = TokenNode
type RegisterParser = TokenParser[RegisterNode]

func NewRegisterParser() Parser[RegisterNode] {
	return RegisterParser{lex.RegisterToken}
}
