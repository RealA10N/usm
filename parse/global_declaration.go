package parse

import "alon.kr/x/usm/source"

// MARK: Node

type GlobalDeclarationNode struct {
	Identifier GlobalNode
	Type       TypeNode
	Immediate  *ImmediateValueNode
}

func (n GlobalDeclarationNode) stringImmediate(ctx *StringContext) string {
	if n.Immediate == nil {
		return ""
	} else {
		return " " + (*n.Immediate).String(ctx)
	}
}

func (n GlobalDeclarationNode) String(ctx *StringContext) string {
	id := n.Identifier.String(ctx)
	typ := n.Type.String(ctx)
	imm := n.stringImmediate(ctx)
	return id + " " + typ + imm
}

func (n GlobalDeclarationNode) View() source.UnmanagedSourceView {
	if n.Immediate == nil {
		return n.Identifier.View().MergeEnd(n.Type.View())
	} else {
		return n.Identifier.View().MergeEnd((*n.Immediate).View())
	}
}
