package parse

import (
	"usm/lex"
	"usm/source"
)

type ArgumentNode struct {
	View     source.UnmanagedSourceView
	Type     lex.Token
	Register lex.Token
}

type ArgumentParser struct{}

func (ArgumentParser) Parse(v *TokenView) (node ArgumentNode, err ParsingError) {
	typ, err := ConsumeToken(v, lex.TypToken)
	if err != nil {
		return
	}

	reg, err := ConsumeToken(v, lex.RegToken)
	if err != nil {
		return
	}

	node = ArgumentNode{
		View:     SourceViewFromBoundaryTokens(typ, reg),
		Type:     typ,
		Register: reg,
	}

	return node, nil
}
