package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type ImmediateNode struct {
	Type  TypeNode
	Value source.UnmanagedSourceView
}

func (n ImmediateNode) View() source.UnmanagedSourceView {
	return n.Type.Merge(n.Value)
}

func (n ImmediateNode) String(ctx source.SourceContext) string {
	return n.Type.String(ctx) + " " + string(n.Value.Raw(ctx))
}

type ImmediateParser struct {
	TypeParser TypeParser
}

func (ImmediateParser) String() string {
	return "immediate ($type #value)"
}

func (p ImmediateParser) Parse(v *TokenView) (node ImmediateNode, err ParsingError) {
	node.Type, err = p.TypeParser.Parse(v)
	if err != nil {
		return
	}

	tkn, err := v.ConsumeToken(lex.ImmediateToken)
	if err != nil {
		return
	}

	node.Value = tkn.View
	return node, nil
}
