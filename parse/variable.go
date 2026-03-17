package parse

import "alon.kr/x/usm/lex"

type VariableNode struct{ TokenNode }
type VariableParser struct{ TokenParser[VariableNode] }

func VariableNodeCreator(tkn lex.Token) VariableNode {
	return VariableNode{TokenNode{tkn.View}}
}

func NewVariableParser() Parser[VariableNode] {
	return VariableParser{
		TokenParser: TokenParser[VariableNode]{
			Token:       lex.VariableToken,
			NodeCreator: VariableNodeCreator,
		},
	}
}
