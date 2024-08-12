package parse

import "alon.kr/x/usm/source"

type TypeFieldNode struct {
	Type   TypeNode
	Labels []LabelNode
}

func (n TypeFieldNode) View() source.UnmanagedSourceView {
	v := n.Type.View()

	if len(n.Labels) > 0 {
		v.Start = n.Labels[0].View().Start
	}

	return v
}

func (n TypeFieldNode) stringLabels(ctx source.SourceContext) (s string) {
	for _, label := range n.Labels {
		s += label.String(ctx) + " "
	}

	return
}

func (n TypeFieldNode) String(ctx source.SourceContext) (s string) {
	return "\t" + n.stringLabels(ctx) + n.Type.String(ctx) + "\n"
}

type TypeFieldParser struct {
	LabelParser LabelParser
	TypeParser  TypeParser
}

func (TypeFieldParser) String() (s string) {
	return "type field"
}

func (p TypeFieldParser) Parse(v *TokenView) (node TypeFieldNode, err ParsingError) {
	node.Labels, _ = ParseManyIgnoreSeparators(p.LabelParser, v)
	node.Type, err = p.TypeParser.Parse(v)
	return
}
