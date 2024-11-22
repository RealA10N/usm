package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

// MARK: Var

type VarNode = TokenNode
type VarParser = TokenParser[VarNode]

func NewVarParser() VarParser {
	return VarParser{lex.VarKeywordToken}
}

// MARK: Declaration

type VarDeclarationNode struct {
	Declaration GlobalDeclarationNode
}

func (n VarDeclarationNode) View() core.UnmanagedSourceView {
	// TODO: not accurate, this does not include the 'var' keyword.
	return n.Declaration.View()
}

func (n VarDeclarationNode) String(ctx *StringContext) string {
	return "var " + n.Declaration.String(ctx)
}

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

	declaration, err := p.GlobalDeclarationParser.Parse(v)
	if err != nil {
		return
	}

	node = VarDeclarationNode{Declaration: declaration}
	return node, nil
}
