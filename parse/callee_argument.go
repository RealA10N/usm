package parse

import (
	"fmt"
	"usm/source"
)

type CalleeArgumentNode struct {
	Type     TypeNode
	Register RegisterNode
}

func (n CalleeArgumentNode) View() source.UnmanagedSourceView {
	return n.Type.View().Merge(n.Register.View())
}

func (n CalleeArgumentNode) String(ctx source.SourceContext) string {
	return fmt.Sprintf("%s %s", n.Type.String(ctx), n.Register.String(ctx))
}

type CalleeArgumentParser struct{}

func (CalleeArgumentParser) Parse(v *TokenView) (node CalleeArgumentNode, err ParsingError) {
	typ, err := TypeParser{}.Parse(v)
	if err != nil {
		return
	}

	reg, err := RegisterParser{}.Parse(v)
	if err != nil {
		return
	}

	node = CalleeArgumentNode{
		Type:     typ,
		Register: reg,
	}

	return
}
