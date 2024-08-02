package parse

import (
	"usm/lex"
	"usm/source"
)

type ArgumentNode struct {
	Node
	Type     lex.Token
	Register lex.Token
}

type ArgumentNodeParser struct{}

func ConsumeToken(v *TokenView, typ lex.TokenType) (tkn lex.Token, perr ParsingError) {
	tknView, restView := v.Partition(1)
	tkn, err := tknView.At(0)

	if err != nil {
		perr = EofError{Expected: lex.TypToken}
		return
	}

	if tkn.Type != typ {
		perr = UnexpectedTokenError{Expected: typ, Got: tkn}
		return
	}

	*v = restView
	return tkn, nil
}

func NewNodeFromBoundaryTokens(first, last lex.Token) Node {
	return Node{
		View: source.UnmanagedSourceView{
			Start: first.View.Start,
			End:   last.View.End},
	}
}

func (p ArgumentNodeParser) Parse(v *TokenView) (node ArgumentNode, err ParsingError) {
	typ, err := ConsumeToken(v, lex.TypToken)
	if err != nil {
		return
	}

	reg, err := ConsumeToken(v, lex.RegToken)
	if err != nil {
		return
	}

	node = ArgumentNode{
		Node:     NewNodeFromBoundaryTokens(typ, reg),
		Type:     typ,
		Register: reg,
	}

	return node, nil
}
