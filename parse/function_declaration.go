package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type FunctionDeclarationNode struct {
	source.UnmanagedSourceView
	Identifier source.UnmanagedSourceView
	Parameters []ParameterNode
	Returns    []TypeNode
}

func (n FunctionDeclarationNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n FunctionDeclarationNode) String(ctx source.SourceContext) (s string) {
	for _, ret := range n.Returns {
		s += ret.String(ctx) + " "
	}

	s += string(n.Identifier.Raw(ctx))

	for _, arg := range n.Parameters {
		s += " " + arg.String(ctx)
	}

	return
}

type FunctionDeclarationParser struct {
	ParameterParser ParameterParser
	TypeParser      TypeParser
}

func (FunctionDeclarationParser) parseIdentifier(v *TokenView, node *FunctionDeclarationNode) ParsingError {
	id, err := v.ConsumeToken(lex.GlobalToken)
	if err != nil {
		return err
	}

	node.Identifier = id.View
	return nil
}

func (FunctionDeclarationParser) updateNodeViewStart(node *FunctionDeclarationNode) {
	if len(node.Returns) > 0 {
		node.Start = node.Returns[0].View().Start
	} else {
		node.Start = node.Identifier.Start
	}
}

func (FunctionDeclarationParser) updateNodeViewEnd(node *FunctionDeclarationNode) {
	if len(node.Parameters) > 0 {
		node.End = node.Parameters[len(node.Parameters)-1].View().End
	} else {
		node.End = node.Identifier.End
	}
}

func (p FunctionDeclarationParser) Parse(v *TokenView) (node FunctionDeclarationNode, err ParsingError) {
	node.Returns = ParseMany(p.TypeParser, v)

	err = p.parseIdentifier(v, &node)
	if err != nil {
		return
	}

	node.Parameters = ParseMany(p.ParameterParser, v)
	p.updateNodeViewStart(&node)
	p.updateNodeViewEnd(&node)
	return
}
