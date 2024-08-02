package parse

import (
	"usm/source"
)

// TODO: add label support

type BlockNode struct {
	source.UnmanagedSourceView
	Instructions []InstructionNode
}

func (n BlockNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n BlockNode) String(ctx source.SourceContext) (s string) {
	s = "{\n"
	for _, inst := range n.Instructions {
		s += "\t" + inst.String(ctx) + "\n"
	}
	s += "}\n"
	return s
}
