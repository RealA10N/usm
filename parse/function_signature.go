package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type FunctionSignatureNode struct {
	source.UnmanagedSourceView
	Identifier source.UnmanagedSourceView
	Parameters []ParameterNode
	Returns    []TypeNode
}

func (n FunctionSignatureNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n FunctionSignatureNode) String(ctx source.SourceContext) (s string) {
	for _, ret := range n.Returns {
		s += ret.String(ctx) + " "
	}

	s += string(n.Identifier.Raw(ctx))

	for _, arg := range n.Parameters {
		s += " " + arg.String(ctx)
	}

	return
}

type FunctionSignatureParser struct {
	ParameterParser ParameterParser
	TypeParser      TypeParser
}

func (FunctionSignatureParser) String() string {
	return "function signature"
}

func (FunctionSignatureParser) parseIdentifier(v *TokenView, node *FunctionSignatureNode) ParsingError {
	id, err := v.ConsumeToken(lex.GlobalToken)
	if err != nil {
		return err
	}

	node.Identifier = id.View
	return nil
}

func (FunctionSignatureParser) updateNodeViewStart(node *FunctionSignatureNode) {
	if len(node.Returns) > 0 {
		node.Start = node.Returns[0].View().Start
	} else {
		node.Start = node.Identifier.Start
	}
}

func (FunctionSignatureParser) updateNodeViewEnd(node *FunctionSignatureNode) {
	if len(node.Parameters) > 0 {
		node.End = node.Parameters[len(node.Parameters)-1].View().End
	} else {
		node.End = node.Identifier.End
	}
}

func (p FunctionSignatureParser) Parse(v *TokenView) (node FunctionSignatureNode, err ParsingError) {
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
