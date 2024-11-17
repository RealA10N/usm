package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

// MARK: Const

type ConstNode = TokenNode

func NewConstParser() Parser[ConstNode] {
	return TokenParser[ConstNode]{lex.ConstKeywordToken}
}

// MARK: Declaration
// TODO: the ConstDeclaration and VarDeclaration are very similar, and there is
// a lot of duplicated code.

type ConstDeclarationNode struct {
	Declaration GlobalDeclarationNode
}

func (n ConstDeclarationNode) View() core.UnmanagedSourceView {
	// TODO: not accurate, this does not include the 'const' keyword.
	return n.Declaration.View()
}

func (n ConstDeclarationNode) String(ctx *StringContext) string {
	return "const " + n.Declaration.String(ctx)
}

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

	declaration, err := p.GlobalDeclarationParser.Parse(v)
	if err != nil {
		return
	}

	node = ConstDeclarationNode{Declaration: declaration}
	return node, nil
}
