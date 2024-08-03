package parse

import (
	"usm/lex"
	"usm/source"
)

type FunctionNode struct {
	source.UnmanagedSourceView
	Signature SignatureNode
	Block     BlockNode
}

func (n FunctionNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n FunctionNode) String(ctx source.SourceContext) string {
	return "def " + n.Signature.String(ctx) + " " + n.Block.String(ctx)
}

type FunctionParser struct {
	SignatureParser SignatureParser
	BlockParser     BlockParser
}

func (FunctionParser) parseDef(v *TokenView, node *FunctionNode) ParsingError {
	def, err := v.ConsumeToken(lex.DefToken)
	if err != nil {
		return err
	}

	node.Start = def.View.Start
	return nil
}

func (p FunctionParser) Parse(v *TokenView) (node FunctionNode, err ParsingError) {
	err = p.parseDef(v, &node)
	if err != nil {
		return
	}

	node.Signature, err = p.SignatureParser.Parse(v)
	if err != nil {
		return
	}

	node.Block, err = p.BlockParser.Parse(v)
	if err != nil {
		return
	}

	node.End = node.Block.View().End
	return
}
