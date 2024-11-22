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

// MARK: Parser

type TypeDecoratorParser struct{}

func NewTypeDecoratorParser() TypeDecoratorParser {
	return TypeDecoratorParser{}
}

func (p TypeDecoratorParser) Parse(v *TokenView) (node TypeDecoratorNode, err ParsingError) {
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
		err = core.GenericError{
			ErrorMessage:  "unexpected type decorator",
			ErrorLocation: tkn.View,
			IsInternal:    true,
		}
		return
	}

	return node, nil
}
