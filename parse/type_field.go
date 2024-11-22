package parse

import (
	"strings"

	"alon.kr/x/usm/core"
)

type TypeFieldNode struct {
	Type   TypeNode
	Labels []LabelNode
}

func (n TypeFieldNode) View() core.UnmanagedSourceView {
	v := n.Type.View()

	if len(n.Labels) > 0 {
		v.Start = n.Labels[0].View().Start
	}

	return v
}

func (n TypeFieldNode) stringLabels(ctx *StringContext) (s string) {
	for _, label := range n.Labels {
		s += label.String(ctx) + " "
	}

	return
}

func (n TypeFieldNode) String(ctx *StringContext) (s string) {
	prefix := strings.Repeat("\t", ctx.Indent)
	labels := n.stringLabels(ctx)
	typ := n.Type.String(ctx)
	return prefix + labels + typ + "\n"
}

type TypeFieldParser struct {
	LabelParser Parser[LabelNode]
	TypeParser  TypeParser
}

func NewTypeFieldParser() Parser[TypeFieldNode] {
	return TypeFieldParser{
		LabelParser: NewLabelParser(),
	}
}

func (p TypeFieldParser) Parse(v *TokenView) (node TypeFieldNode, err core.Result) {
	node.Labels, _ = ParseManyIgnoreSeparators(p.LabelParser, v)
	node.Type, err = p.TypeParser.Parse(v)
	return
}
