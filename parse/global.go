package parse

import (
	"alon.kr/x/usm/lex"
)

type GlobalNode = TokenNode
type GlobalParser = TokenParser[GlobalNode]

func NewGlobalParser() GlobalParser {
	return GlobalParser{lex.GlobalToken}
}
