package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

// MARK: Var

type VarNode struct{ TokenNode }
type VarParser struct{ TokenParser[VarNode] }

func VarNodeCreator(tkn lex.Token) VarNode {
	return VarNode{TokenNode{tkn.View}}
}

func NewVarParser() Parser[VarNode] {
	return VarParser{
		TokenParser[VarNode]{
			Token:       lex.VarKeywordToken,
			NodeCreator: VarNodeCreator,
		},
	}
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
	err core.Result,
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
