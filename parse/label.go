package parse

import (
	"alon.kr/x/usm/lex"
)

type LabelNode = TokenNode

func NewLabelParser() Parser[LabelNode] {
	return TokenParser[LabelNode]{lex.LabelToken}
}
