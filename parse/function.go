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

type FunctionParser struct {
	SignatureParser SignatureParser
	BlockParser     BlockParser
}

func (p FunctionParser) Parse(v *TokenView) (node FunctionNode, err ParsingError) {
	node.Signature, err = p.SignatureParser.Parse(v)
	if err != nil {
		return
	}

	node.Block, err = p.BlockParser.Parse(v)
	return
}
