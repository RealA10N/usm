package parse

import "alon.kr/x/usm/core"

type TargetNode struct {
	// Optional type declaration. Depending on the instruction, the type may be
	// inferred and does not need to be provided explicitly.
	Type *TypeNode

	Register RegisterNode
}

func (n TargetNode) View() core.UnmanagedSourceView {
	v := n.Register.View()
	if n.Type != nil {
		v = v.MergeStart(n.Type.View())
	}
	return v
}

func (n TargetNode) String(ctx *StringContext) (s string) {
	if n.Type != nil {
		s = n.Type.String(ctx) + " "
	}
	return s + n.Register.String(ctx)
}
