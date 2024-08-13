package parse

import (
	"alon.kr/x/usm/lex"
)

type GlobalNode = TokenNode

func NewGlobalParser() Parser[GlobalNode] {
	return TokenParser[GlobalNode]{lex.GlobalToken}
}
