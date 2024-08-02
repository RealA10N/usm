package parse

import "usm/source"

type SignatureNode struct {
	source.UnmanagedSourceView
	Identifier source.UnmanagedSourceView
	Arguments  []CalleeArgumentNode
	Returns    []TypeNode
}

func (n SignatureNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n SignatureNode) String(ctx source.SourceContext) string {
	s := "def "
	for _, ret := range n.Returns {
		s += ret.String(ctx) + " "
	}

	s += string(n.Identifier.Raw(ctx))

	for _, arg := range n.Arguments {
		s += " " + arg.String(ctx)
	}

	return s
}
