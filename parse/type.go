package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

// MARK: Node

type TypeNode struct {
	Identifier core.UnmanagedSourceView
	Decorators []TypeDecoratorNode
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
	for _, dec := range n.Decorators {
		s += " " + dec.String(ctx)
	}
	return s
}

// MARK: Parser

type TypeParser struct {
	TypeDecoratorParser
}

func NewTypeParser() TypeParser {
	return TypeParser{
		TypeDecoratorParser: NewTypeDecoratorParser(),
	}
}

func (p TypeParser) Parse(v *TokenView) (node TypeNode, err ParsingError) {
	tkn, err := v.ConsumeToken(lex.TypeToken)
	if err != nil {
		return
	}

	node.Identifier = tkn.View

	for err == nil {
		// TODO: there is not distinction here between an error in the decorator parsing
		// and the end of available decorators. In particular, any decorator parsing error
		// will get swallowed here and not shown to the user.
		var decorator TypeDecoratorNode
		decorator, err = p.TypeDecoratorParser.Parse(v)
		if err == nil {
			node.Decorators = append(node.Decorators, decorator)
		}
	}

	return node, nil
}
