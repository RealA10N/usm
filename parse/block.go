package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type BlockNode[NodeT Node] struct {
	core.UnmanagedSourceView
	Nodes []NodeT
	// TrailingComments holds whole-line comments before the closing '}'.
	TrailingComments []lex.Comment
}

func (n BlockNode[NodeT]) View() core.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n BlockNode[NodeT]) String(ctx *StringContext) (s string) {
	if len(n.Nodes) == 0 && len(n.TrailingComments) == 0 {
		return "{ }"
	}

	s = "{\n"
	ctx.Indent++
	for _, node := range n.Nodes {
		s += node.String(ctx)
	}
	s += ctx.renderComments(n.TrailingComments)
	ctx.Indent--
	s += ctx.indent() + "}"
	return
}

// commentAttachable is implemented by node types that store leading comments.
// BlockParser uses this interface to attach collected comments to parsed nodes.
type commentAttachable interface {
	attachLeadingComments([]lex.Comment)
}

type BlockParser[NodeT Node] struct {
	Parser Parser[NodeT]
}

// parseBlockNodes parses all nodes inside a block, collecting whole-line
// comments as leading comments on the following node. Comments that appear
// after the last node (before '}') are returned as trailing block comments.
//
// The node type may optionally implement commentAttachable; if it does,
// leading comments are attached directly to each node. Otherwise they are
// silently dropped (acceptable for blocks like type fields or immediates).
func (p BlockParser[NodeT]) parseBlockNodes(v *TokenView) (nodes []NodeT, trailing []lex.Comment) {
	for {
		pending := v.consumeLeadingComments()

		// Peek: if the next token is '}' or the view is empty, the pending
		// comments are trailing block comments, not leading instruction comments.
		front, err := v.At(0)
		if err != nil || front.Type == lex.RightCurlyBraceToken {
			trailing = pending
			return
		}

		node, parseErr := p.Parser.Parse(v)
		if parseErr != nil {
			trailing = pending
			return
		}

		if ca, ok := any(&node).(commentAttachable); ok {
			ca.attachLeadingComments(pending)
		}
		nodes = append(nodes, node)
	}
}

func (p BlockParser[NodeT]) Parse(v *TokenView) (block BlockNode[NodeT], err core.Result) {
	leftCurly, err := v.ConsumeToken(lex.LeftCurlyBraceToken)
	if err != nil {
		return
	}

	block.Start = leftCurly.View.Start
	block.Nodes, block.TrailingComments = p.parseBlockNodes(v)

	rightCurly, err := v.ConsumeTokenIgnoreSeparator(lex.RightCurlyBraceToken)
	if err != nil {
		return
	}

	block.End = rightCurly.View.End
	return block, nil
}
