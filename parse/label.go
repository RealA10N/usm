package parse

import (
	"alon.kr/x/usm/lex"
)

type LabelNode struct{ TokenNode }
type LabelParser struct{ TokenParser[LabelNode] }

func LabelNodeCreator(tkn lex.Token) LabelNode {
	return LabelNode{TokenNode{tkn.View}}
}

func NewLabelParser() Parser[LabelNode] {
	return LabelParser{
		TokenParser: TokenParser[LabelNode]{
			Token:       lex.LabelToken,
			NodeCreator: LabelNodeCreator,
		},
	}
}
