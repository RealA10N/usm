package parse

import (
	"usm/lex"
	"usm/source"
)

// TODO: this should be a tagged union.
// The supported caller arguments are registers, globals, immediates (and labales?)

type CallerArgumentNode struct {
	source.UnmanagedSourceView
}

func (n CallerArgumentNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n CallerArgumentNode) String(ctx source.SourceContext) string {
	return string(n.UnmanagedSourceView.Raw(ctx))
}

type CallerArgumentParser struct{}

func (CallerArgumentParser) Parse(v *TokenView) (node CallerArgumentNode, err ParsingError) {
	tkn, err := ConsumeToken(v, lex.RegToken, lex.ImmToken, lex.GlbToken)
	if err != nil {
		return
	}

	node = CallerArgumentNode{tkn.View}
	return
}
