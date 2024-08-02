package parse

import (
	"usm/source"
)

type InstructionNode struct {
	Operation source.UnmanagedSourceView
	Arguments []CallerArgumentNode
	Targets   []RegisterNode
}

func (n InstructionNode) View() source.UnmanagedSourceView {
	first := n.Operation
	last := n.Operation

	if len(n.Targets) > 0 {
		first = n.Targets[0].View()
	}

	if len(n.Arguments) > 0 {
		last = n.Arguments[len(n.Arguments)-1].View()
	}

	return first.Merge(last)
}

func (n InstructionNode) stringArguments(ctx source.SourceContext) (s string) {
	if len(n.Arguments) == 0 {
		return
	}

	for _, arg := range n.Arguments {
		s += " " + arg.String(ctx)
	}

	return
}

func (n InstructionNode) stringTargets(ctx source.SourceContext) (s string) {
	if len(n.Targets) == 0 {
		return
	}

	for _, tgt := range n.Targets {
		s += tgt.String(ctx) + " "
	}

	s += "= "
	return
}

func (n InstructionNode) String(ctx source.SourceContext) string {
	op := string(n.Operation.Raw(ctx))
	return n.stringTargets(ctx) + op + n.stringArguments(ctx)
}
