package parse

import "alon.kr/x/usm/source"

// MARK: Node

type GlobalDeclarationNode struct {
	Identifier GlobalNode
	Type       TypeNode
	Immediate  *ImmediateValueNode
}

func (n GlobalDeclarationNode) stringImmediate(ctx *StringContext) string {
	if n.Immediate == nil {
		return ""
	} else {
		return " " + (*n.Immediate).String(ctx)
	}
}

func (n GlobalDeclarationNode) String(ctx *StringContext) string {
	id := n.Identifier.String(ctx)
	typ := n.Type.String(ctx)
	imm := n.stringImmediate(ctx)
	return id + " " + typ + imm
}

func (n GlobalDeclarationNode) View() source.UnmanagedSourceView {
	if n.Immediate == nil {
		return n.Identifier.View().MergeEnd(n.Type.View())
	} else {
		return n.Identifier.View().MergeEnd((*n.Immediate).View())
	}
}

// MARK: Parser

type GlobalDeclarationParser struct {
	GlobalParser         GlobalParser
	TypeParser           TypeParser
	ImmediateValueParser *ImmediateValueParser
}

func NewGlobalDeclarationParser() GlobalDeclarationParser {
	return GlobalDeclarationParser{
		ImmediateValueParser: NewImmediateValueParser(),
	}
}

func (p *GlobalDeclarationParser) Parse(v *TokenView) (
	node GlobalDeclarationNode,
	err ParsingError,
) {
	node.Identifier, err = p.GlobalParser.Parse(v)
	if err != nil {
		return
	}

	node.Type, err = p.TypeParser.Parse(v)
	if err != nil {
		return
	}

	immediate, err := p.ImmediateValueParser.Parse(v)
	if err == nil {
		// TODO: improve error catching. Not every error here means that there
		// is no immediate value! what if an immediate block is provided (the
		// next token is '{') but there is an error parsing the ImmediateBlock?
		node.Immediate = &immediate
	}

	return node, nil
}
