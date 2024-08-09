package parse

import (
	"fmt"

	"alon.kr/x/usm/source"
)

type ParameterNode struct {
	Type     TypeNode
	Register RegisterNode
}

func (n ParameterNode) View() source.UnmanagedSourceView {
	return n.Type.View().Merge(n.Register.View())
}

func (n ParameterNode) String(ctx source.SourceContext) string {
	return fmt.Sprintf("%s %s", n.Type.String(ctx), n.Register.String(ctx))
}

type ParameterParser struct {
	TypeParser     TypeParser
	RegisterParser RegisterParser
}

func (ParameterParser) String() string {
	return "parameter ($type %register)"
}

func (p ParameterParser) Parse(v *TokenView) (node ParameterNode, err ParsingError) {
	typ, err := p.TypeParser.Parse(v)
	if err != nil {
		return
	}

	reg, err := p.RegisterParser.Parse(v)
	if err != nil {
		return
	}

	return ParameterNode{typ, reg}, nil
}
