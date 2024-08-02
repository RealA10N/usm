package parse

import (
	"usm/source"
)

type FunctionNode struct {
	Signature SignatureNode
	Block     BlockNode
}

func (n FunctionNode) View() source.UnmanagedSourceView {
	return n.Signature.View().Merge(n.Block.View())
}

func (n FunctionNode) String(ctx source.SourceContext) string {
	return n.Signature.String(ctx) + " " + n.Block.String(ctx)
}
