package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type ArgumentNode struct {
	source.UnmanagedSourceView
}

func (n ArgumentNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n ArgumentNode) String(ctx source.SourceContext) string {
	return string(n.UnmanagedSourceView.Raw(ctx))
}

type ArgumentParser struct{}

// TODO: fix argument to use appropriate subparsers
// TODO: immediate argument should be prefixed with a type

func (ArgumentParser) Parse(v *TokenView) (node ArgumentNode, err ParsingError) {
	tkn, err := v.ConsumeToken(lex.RegisterToken, lex.ImmediateToken, lex.GlobalToken)
	if err != nil {
		return
	}

	node = ArgumentNode{tkn.View}
	return
}
