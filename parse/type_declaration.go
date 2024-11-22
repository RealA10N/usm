package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

// MARK: Node

type TypeDeclarationNode struct {
	core.UnmanagedSourceView
	Identifier core.UnmanagedSourceView
	Fields     BlockNode[TypeFieldNode]
}

func (n TypeDeclarationNode) View() core.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n TypeDeclarationNode) String(ctx *StringContext) string {
	id := string(n.Identifier.Raw(ctx.SourceContext))
	fields := n.Fields.String(ctx)
	return "type " + id + " " + fields
}

// MARK: Parser

type TypeDeclarationParser struct {
	FieldsParser BlockParser[TypeFieldNode]
}

func NewTypeDeclarationParser() TypeDeclarationParser {
	return TypeDeclarationParser{
		FieldsParser: BlockParser[TypeFieldNode]{
			Parser: NewTypeFieldParser(),
		},
	}
}

func (TypeDeclarationParser) parseTypeKeyword(v *TokenView, node *TypeDeclarationNode) (err core.Result) {
	kw, err := v.ConsumeToken(lex.TypeKeywordToken)
	if err != nil {
		return
	}

	node.Start = kw.View.Start
	return
}

func (TypeDeclarationParser) parseIdentifier(v *TokenView, node *TypeDeclarationNode) (err core.Result) {
	id, err := v.ConsumeToken(lex.TypeToken)
	if err != nil {
		return
	}

	node.Identifier = id.View
	return
}

func (p TypeDeclarationParser) parseBlock(v *TokenView, node *TypeDeclarationNode) (err core.Result) {
	node.Fields, err = p.FieldsParser.Parse(v)
	if err != nil {
		return
	}

	node.End = node.Fields.End
	return
}

func (p TypeDeclarationParser) Parse(v *TokenView) (node TypeDeclarationNode, err core.Result) {
	err = p.parseTypeKeyword(v, &node)
	if err != nil {
		return
	}

	err = p.parseIdentifier(v, &node)
	if err != nil {
		return
	}

	err = p.parseBlock(v, &node)
	return
}
