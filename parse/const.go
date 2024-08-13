package parse

import "alon.kr/x/usm/lex"

// MARK: Const

type ConstNode = TokenNode

func NewConstParser() Parser[ConstNode] {
	return TokenParser[ConstNode]{lex.ConstKeywordToken}
}

// MARK: Declaration

type ConstDeclarationNode = GlobalDeclarationNode

type ConstDeclarationParser struct {
	ConstParser             Parser[ConstNode]
	GlobalDeclarationParser GlobalDeclarationParser
}

func NewConstDeclarationParser() Parser[ConstDeclarationNode] {
	return ConstDeclarationParser{
		ConstParser:             NewConstParser(),
		GlobalDeclarationParser: NewGlobalDeclarationParser(),
	}
}

func (p ConstDeclarationParser) Parse(v *TokenView) (
	node ConstDeclarationNode,
	err ParsingError,
) {
	_, err = p.ConstParser.Parse(v)
	if err != nil {
		return
	}

	global, err := p.GlobalDeclarationParser.Parse(v)
	if err != nil {
		return
	}

	node = ConstDeclarationNode(global)
	return node, nil
}
