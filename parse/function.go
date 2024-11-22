package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type FunctionNode struct {
	core.UnmanagedSourceView
	Signature    FunctionSignatureNode
	Instructions *BlockNode[InstructionNode]
}

func (n FunctionNode) View() core.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n FunctionNode) String(ctx *StringContext) string {
	s := "func " + n.Signature.String(ctx)
	if n.Instructions != nil {
		s += " " + n.Instructions.String(ctx)
	}

	return s
}

type FunctionParser struct {
	FunctionSignatureParser FunctionSignatureParser
	InstructionBlockParser  BlockParser[InstructionNode]
}

func NewFunctionParser() FunctionParser {
	return FunctionParser{
		FunctionSignatureParser: NewFunctionSignatureParser(),
		InstructionBlockParser: BlockParser[InstructionNode]{
			Parser: NewInstructionParser(),
		},
	}
}

func (FunctionParser) parseFunctionKeyword(v *TokenView, node *FunctionNode) core.Result {
	kw, err := v.ConsumeToken(lex.FuncKeywordToken)
	if err != nil {
		return err
	}

	node.Start = kw.View.Start
	return nil
}

func (p FunctionParser) parseBlock(v *TokenView, node *FunctionNode) {
	instructions, err := p.InstructionBlockParser.Parse(v)
	if err == nil {
		node.Instructions = &instructions
		node.End = node.Instructions.View().End
	} else {
		node.End = node.Signature.View().End
	}
}

func (p FunctionParser) Parse(v *TokenView) (node FunctionNode, err core.Result) {
	err = p.parseFunctionKeyword(v, &node)
	if err != nil {
		return
	}

	node.Signature, err = p.FunctionSignatureParser.Parse(v)
	if err != nil {
		return
	}

	p.parseBlock(v, &node)
	return
}
