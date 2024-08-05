package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

// TODO: this should be a tagged union.
// The supported caller arguments are registers, globals, immediates (and labels?)

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

func (ArgumentParser) Parse(v *TokenView) (node ArgumentNode, err ParsingError) {
	tkn, err := v.ConsumeToken(lex.RegToken, lex.ImmToken, lex.GlbToken)
	if err != nil {
		return
	}

	node = ArgumentNode{tkn.View}
	return
}
