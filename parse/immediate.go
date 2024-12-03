// The structure and syntax of a valid immediate custom type initialization is
// quite complicated and is defined recursively. The following are examples of
// a valid immediate, with the corresponding AST:

//                     ┌---┬--> ImmediateValueNode
// const @constant $32 #1337
//                 └-------┴--> ImmediateNode
// └-----------------------┴--> ConstantNode

// const @global $outer {                     ----┐
//   ┌----------┬--> ImmediateFieldNode           ├ ImmediateBlockNode (not
//   .value #1234                                 | including $outer type),
//          └---┴--> ImmediateFieldValueNode      | and ImmediateNode including
//   .inner { .value #0 }                         | the $outer prefix.
//          └-----------┴--> ImmediateBlockNode   |
// }                                          ----┘

package parse

import (
	"strings"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

// MARK: Final Value
// ImmediateFinalValue Node & Parser are responsible for the #immediate token only.

type ImmediateFinalValueNode struct {
	core.UnmanagedSourceView
}

func (n ImmediateFinalValueNode) View() core.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n ImmediateFinalValueNode) String(ctx *StringContext) string {
	return string(n.UnmanagedSourceView.Raw(ctx.SourceContext))
}

type ImmediateFinalValueParser struct{}

func (ImmediateFinalValueParser) Parse(v *TokenView) (
	node ImmediateFinalValueNode,
	err core.Result,
) {
	tkn, err := v.ConsumeToken(lex.ImmediateToken)
	if err != nil {
		return
	}

	return ImmediateFinalValueNode{tkn.View}, nil
}

// MARK: Value

// This is an interface of the type that appear as a value in a field of a custom
// type initialization. It can be either
// (1.) an ImmediateFinalValueNode (#1234), or
// (2.) an ImmediateBlockNode ({ ... }).
type ImmediateValueNode interface {
	Node
}

type ImmediateValueParser struct {
	ImmediateFinalValueParser *ImmediateFinalValueParser
	ImmediateBlockParser      *ImmediateBlockParser
}

func (p ImmediateValueParser) Parse(v *TokenView) (
	node ImmediateValueNode,
	err core.Result,
) {
	tkn, err := v.PeekToken(lex.ImmediateToken, lex.LeftCurlyBraceToken)
	if err != nil {
		// TODO: improve propagated error message
		return
	}

	switch tkn.Type {
	case lex.ImmediateToken:
		return p.ImmediateFinalValueParser.Parse(v)
	case lex.LeftCurlyBraceToken:
		return p.ImmediateBlockParser.Parse(v)
	default:
		panic("unreachable")
	}
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

func (n ImmediateFieldNode) View() core.UnmanagedSourceView {
	if n.Label != nil {
		return n.Label.View().MergeEnd(n.Value.View())
	} else {
		return n.Value.View()
	}
}

func (n ImmediateFieldNode) stringLabel(ctx *StringContext) (s string) {
	if n.Label != nil {
		return n.Label.String(ctx) + " "
	}
	return
}

func (n ImmediateFieldNode) String(ctx *StringContext) string {
	prefix := strings.Repeat("\t", ctx.Indent)
	label := n.stringLabel(ctx)
	value := n.Value.String(ctx)
	return prefix + label + value + "\n"
}

type ImmediateFieldParser struct {
	LabelParser          Parser[LabelNode]
	ImmediateValueParser *ImmediateValueParser
}

func (p ImmediateFieldParser) parseLabel(v *TokenView, node *ImmediateFieldNode) {
	label, err := p.LabelParser.Parse(v)
	if err != nil {
		return
	}

	node.Label = &label
}

func (p ImmediateFieldParser) Parse(v *TokenView) (node ImmediateFieldNode, err core.Result) {
	p.parseLabel(v, &node)

	val, err := p.ImmediateValueParser.Parse(v)
	if err != nil {
		return
	}

	node.Value = val
	return
}

// MARK: Block

type ImmediateBlockNode = BlockNode[ImmediateFieldNode]
type ImmediateBlockParser = BlockParser[ImmediateFieldNode]

// MARK: Immediate

type ImmediateNode struct {
	Type  TypeNode
	Value ImmediateValueNode
}

func (n ImmediateNode) View() core.UnmanagedSourceView {
	return n.Type.View().MergeEnd(n.Value.View())
}

func (n ImmediateNode) String(ctx *StringContext) string {
	return n.Type.String(ctx) + " " + n.Value.String(ctx)
}

type ImmediateParser struct {
	TypeParser
	ImmediateValueParser *ImmediateValueParser
}

func (p ImmediateParser) Parse(v *TokenView) (node ImmediateNode, err core.Result) {
	node.Type, err = p.TypeParser.Parse(v)
	if err != nil {
		return
	}

	node.Value, err = p.ImmediateValueParser.Parse(v)
	if err != nil {
		return
	}

	return node, nil
}

// MARK: New

func NewImmediateParser() *ImmediateParser {
	return &ImmediateParser{
		ImmediateValueParser: NewImmediateValueParser(),
	}
}

func NewImmediateValueParser() *ImmediateValueParser {
	valueParser := &ImmediateValueParser{
		ImmediateFinalValueParser: &ImmediateFinalValueParser{},
	}
	valueParser.ImmediateBlockParser = NewImmediateBlockParser(valueParser)
	return valueParser
}

func NewImmediateBlockParser(valueParser *ImmediateValueParser) *ImmediateBlockParser {
	return &ImmediateBlockParser{NewImmediateFieldParser(valueParser)}
}

func NewImmediateFieldParser(valueParser *ImmediateValueParser) *ImmediateFieldParser {
	return &ImmediateFieldParser{
		LabelParser:          NewLabelParser(),
		ImmediateValueParser: valueParser,
	}
}
