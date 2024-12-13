package parse

import (
	"alon.kr/x/usm/lex"
)

type GlobalNode struct{ TokenNode }
type GlobalParser struct{ TokenParser[GlobalNode] }

func GlobalNodeCreator(tkn lex.Token) GlobalNode {
	return GlobalNode{TokenNode{tkn.View}}
}

func NewGlobalParser() Parser[GlobalNode] {
	return GlobalParser{
		TokenParser: TokenParser[GlobalNode]{
			Token:       lex.GlobalToken,
			NodeCreator: GlobalNodeCreator,
		},
	}
}
