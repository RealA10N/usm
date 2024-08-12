package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

// MARK: Node

type TypeDeclarationNode struct {
	source.UnmanagedSourceView
	Identifier source.UnmanagedSourceView
	Fields     BlockNode[TypeFieldNode]
}

func (n TypeDeclarationNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n TypeDeclarationNode) String(ctx source.SourceContext) string {
	return "type " + string(n.Identifier.Raw(ctx)) + " " + n.Fields.String(ctx)
}

// MARK: Parser

type TypeDeclarationParser struct {
	FieldsParser BlockParser[TypeFieldNode]
}

func NewTypeDeclarationParser() TypeDeclarationParser {
	return TypeDeclarationParser{
		FieldsParser: BlockParser[TypeFieldNode]{
			Parser: TypeFieldParser{},
		},
	}
}

func (TypeDeclarationParser) String() string {
	return "type declaration"
}

func (TypeDeclarationParser) parseTypeKeyword(v *TokenView, node *TypeDeclarationNode) (err ParsingError) {
	kw, err := v.ConsumeToken(lex.TypeKeywordToken)
	if err != nil {
		return
	}

	node.Start = kw.View.Start
	return
}

func (p TypeDeclarationParser) parseIdentifier(v *TokenView, node *TypeDeclarationNode) (err ParsingError) {
	id, err := v.ConsumeToken(lex.TypeToken)
	if err != nil {
		return
	}

	node.Identifier = id.View
	return
}

func (p TypeDeclarationParser) parseBlock(v *TokenView, node *TypeDeclarationNode) (err ParsingError) {
	node.Fields, err = p.FieldsParser.Parse(v)
	if err != nil {
		return
	}

	node.End = node.Fields.End
	return
}

func (p TypeDeclarationParser) Parse(v *TokenView) (node TypeDeclarationNode, err ParsingError) {
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
