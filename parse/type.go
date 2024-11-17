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

type TypeDecorator struct {
	core.UnmanagedSourceView
	Type TypeDecoratorType
}

type TypeNode struct {
	Identifier core.UnmanagedSourceView
	Decorators []TypeDecorator
}

func (n TypeNode) View() core.UnmanagedSourceView {
	if len(n.Decorators) == 0 {
		return n.Identifier
	} else {
		return n.Identifier.MergeEnd(
			n.Decorators[len(n.Decorators)-1].UnmanagedSourceView,
		)
	}
}

func (n TypeNode) String(ctx *StringContext) string {
	s := string(n.Identifier.Raw(ctx.SourceContext))
	for _, decorator := range n.Decorators {
		s += " " + string(decorator.Raw(ctx.SourceContext))
	}
	return s
}

// MARK: Parser

type TypeParser struct{}

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
		// TODO: replace with error message (perhaps internal error?)
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
