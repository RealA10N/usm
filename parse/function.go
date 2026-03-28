package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type FunctionNode struct {
	core.UnmanagedSourceView
	Signature       FunctionSignatureNode
	Instructions    *BlockNode[InstructionNode]
	LeadingComments []lex.Comment
}

func (n FunctionNode) View() core.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n *FunctionNode) attachLeadingComments(c []lex.Comment) {
	n.LeadingComments = c
}

func (n FunctionNode) stringBlock(ctx *StringContext) string {
	hasInstructions := n.Instructions != nil && (len(n.Instructions.Nodes) > 0 || len(n.Instructions.TrailingComments) > 0)

	if !hasInstructions {
		return "{ }"
	}

	s := "{\n"
	ctx.Indent++
	if n.Instructions != nil {
		for _, instr := range n.Instructions.Nodes {
			s += instr.String(ctx)
		}
		s += ctx.renderComments(n.Instructions.TrailingComments)
	}
	ctx.Indent--
	s += ctx.indent() + "}"
	return s
}

func (n FunctionNode) String(ctx *StringContext) string {
	s := ctx.renderComments(n.LeadingComments)
	s += "func " + n.Signature.String(ctx)
	if n.Instructions != nil {
		s += " " + n.stringBlock(ctx)
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
		InstructionBlockParser:  BlockParser[InstructionNode]{Parser: NewInstructionParser()},
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
	leftCurly, err := v.ConsumeToken(lex.LeftCurlyBraceToken)
	if err != nil {
		node.End = node.Signature.View().End
		return
	}

	nodes, trailing := p.InstructionBlockParser.parseBlockNodes(v)

	rightCurly, err := v.ConsumeTokenIgnoreSeparator(lex.RightCurlyBraceToken)
	if err != nil {
		node.End = leftCurly.View.End
		return
	}

	block := &BlockNode[InstructionNode]{
		UnmanagedSourceView: core.UnmanagedSourceView{
			Start: leftCurly.View.Start,
			End:   rightCurly.View.End,
		},
		Nodes:            nodes,
		TrailingComments: trailing,
	}
	node.Instructions = block
	node.End = rightCurly.View.End
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
