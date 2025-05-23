package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

// MARK: Node

type TypeDecoratorType uint8

const (
	PointerTypeDecorator TypeDecoratorType = iota
	RepeatTypeDecorator
)

type TypeDecoratorNode struct {
	core.UnmanagedSourceView
	Type TypeDecoratorType
}

func (n TypeDecoratorNode) String(ctx *StringContext) string {
	return string(n.Raw(ctx.SourceContext))
}

func (n TypeDecoratorNode) View() core.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

// MARK: Parser

type TypeDecoratorParser struct{}

func NewTypeDecoratorParser() TypeDecoratorParser {
	return TypeDecoratorParser{}
}

func (p TypeDecoratorParser) Parse(v *TokenView) (node TypeDecoratorNode, err core.Result) {
	tkn, err := v.ConsumeToken(lex.PointerToken, lex.RepeatToken)
	if err != nil {
		return
	}

	node.UnmanagedSourceView = tkn.View

	switch tkn.Type {
	case lex.PointerToken:
		node.Type = PointerTypeDecorator
	case lex.RepeatToken:
		node.Type = RepeatTypeDecorator
	default:
		err = core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Unexpected type decorator",
			Location: &tkn.View,
		}}
		return
	}

	return node, nil
}
