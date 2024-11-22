package parse

import (
	"alon.kr/x/usm/lex"
)

type LabelNode = TokenNode
type LabelParser = TokenParser[LabelNode]

func NewLabelParser() LabelParser {
	return LabelParser{lex.LabelToken}
}
