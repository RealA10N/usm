package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

// MARK: Value
// ImmediateValue Node & Parser are responsible for the #immediate token only.

type ImmediateValueNode struct {
	source.UnmanagedSourceView
}

func (n ImmediateValueNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n ImmediateValueNode) String(ctx source.SourceContext) string {
	return string(n.UnmanagedSourceView.Raw(ctx))
}

type ImmediateValueParser struct{}

func (ImmediateValueParser) Parse(v *TokenView) (node ImmediateValueNode, err ParsingError) {
	tkn, err := v.ConsumeToken(lex.ImmediateToken)
	if err != nil {
		return
	}

	return ImmediateValueNode{tkn.View}, nil
}

// MARK: Field

// ImmediateFieldNode represents a single field (entry) in the initialization of
// an immediate custom type (struct).
type ImmediateFieldNode struct {
	// At most one (possible, zero) field labels can be specified for each filed
	// Field is nil if a label is not specified.
	Label *LabelNode
	Value ImmediateValueNode
}

func (n ImmediateFieldNode) View() source.UnmanagedSourceView {
	if n.Label != nil {
		return n.Label.View().MergeEnd(n.Value.View())
	} else {
		return n.Value.View()
	}
}

func (n ImmediateFieldNode) String(ctx source.SourceContext) string {
	if n.Label != nil {
		return n.Label.String(ctx) + " " + n.Value.String(ctx)
	} else {
		return n.Value.String(ctx)
	}
}

type ImmediateFieldParser struct {
	LabelParser          LabelParser
	ImmediateValueParser ImmediateValueParser
}

func (p ImmediateFieldParser) tryParsingLabel(v *TokenView, node *ImmediateFieldNode) {
	label, err := p.LabelParser.Parse(v)
	if err != nil {
		return
	}

	node.Label = &label
}

func (p ImmediateFieldParser) Parse(v *TokenView) (node ImmediateFieldNode, err ParsingError) {
	p.tryParsingLabel(v, &node)

	val, err := p.ImmediateValueParser.Parse(v)
	if err != nil {
		return
	}

	node.Value = val
	return
}

// MARK: Immediate

type ImmediateNode struct {
	Type  TypeNode
	Value source.UnmanagedSourceView
}

func (n ImmediateNode) View() source.UnmanagedSourceView {
	return n.Type.View().MergeEnd(n.Value)
}

func (n ImmediateNode) String(ctx source.SourceContext) string {
	return n.Type.String(ctx) + " " + string(n.Value.Raw(ctx))
}

type ImmediateParser struct {
	TypeParser TypeParser
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

// func (p ImmediateFieldParser) Parse(v *TokenView) (node ImmediateNode, err ParsingError) {
// }
