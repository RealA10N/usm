package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

// MARK: Node

type TypeDecoratorType uint8

const (
	PointerTypeDecorator TypeDecoratorType = iota
	RepeatTypeDecorator
)

type TypeDecorator struct {
	source.UnmanagedSourceView
	Type TypeDecoratorType
}

type TypeNode struct {
	Identifier source.UnmanagedSourceView
	Decorators []TypeDecorator
}

func (n TypeNode) View() source.UnmanagedSourceView {
	return n.Identifier
}

func (n TypeNode) String(ctx source.SourceContext) string {
	return string(n.Identifier.Raw(ctx))
}

// MARK: Parser

type TypeParser struct{}

func (TypeParser) String() string {
	return "$type"
}

func (p TypeParser) parseDecorator(v *TokenView, node *TypeNode) (err ParsingError) {
	tkn, err := v.ConsumeToken(lex.PointerToken, lex.RepeatToken)
	if err != nil {
		return
	}

	decorator := TypeDecorator{UnmanagedSourceView: tkn.View}
	switch tkn.Type {
	case lex.PointerToken:
		decorator.Type = PointerTypeDecorator
	case lex.RepeatToken:
		decorator.Type = RepeatTypeDecorator
	default:
		panic("unreachable")
	}

	node.Decorators = append(node.Decorators, decorator)
	return nil
}

func (p TypeParser) Parse(v *TokenView) (node TypeNode, err ParsingError) {
	tkn, err := v.ConsumeToken(lex.TypeToken)
	if err != nil {
		return
	}

	node.Identifier = tkn.View
	for err == nil {
		err = p.parseDecorator(v, &node)
	}

	return node, nil
}
