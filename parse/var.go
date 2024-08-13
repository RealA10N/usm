package parse

import "alon.kr/x/usm/lex"

// MARK: Var

type VarNode = TokenNode

func NewVarParser() Parser[VarNode] {
	return TokenParser[VarNode]{lex.VarKeywordToken}
}

// MARK: Declaration

type VarDeclarationNode = GlobalDeclarationNode

type VarDeclarationParser struct {
	VarParser               Parser[VarNode]
	GlobalDeclarationParser GlobalDeclarationParser
}

func NewVarDeclarationParser() Parser[VarDeclarationNode] {
	return VarDeclarationParser{
		VarParser:               NewVarParser(),
		GlobalDeclarationParser: NewGlobalDeclarationParser(),
	}
}

func (p VarDeclarationParser) Parse(v *TokenView) (
	node VarDeclarationNode,
	err ParsingError,
) {
	_, err = p.VarParser.Parse(v)
	if err != nil {
		return
	}

	global, err := p.GlobalDeclarationParser.Parse(v)
	if err != nil {
		return
	}

	node = VarDeclarationNode(global)
	return node, nil
}
