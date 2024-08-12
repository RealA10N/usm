package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type FunctionNode struct {
	source.UnmanagedSourceView
	Declaration  FunctionDeclarationNode
	Instructions *BlockNode[InstructionNode]
}

func (n FunctionNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n FunctionNode) String(ctx source.SourceContext) string {
	s := "func " + n.Declaration.String(ctx)
	if n.Instructions != nil {
		s += " " + n.Instructions.String(ctx)
	} else {
		s += "\n"
	}

	return s
}

type FunctionParser struct {
	FunctionDeclarationParser FunctionDeclarationParser
	InstructionBlockParser    BlockParser[InstructionNode]
}

func NewFunctionParser() FunctionParser {
	return FunctionParser{
		FunctionDeclarationParser: FunctionDeclarationParser{},
		InstructionBlockParser: BlockParser[InstructionNode]{
			Parser: InstructionParser{},
		},
	}
}

func (FunctionParser) String() string {
	return "function"
}

func (FunctionParser) parseFunctionKeyword(v *TokenView, node *FunctionNode) ParsingError {
	kw, err := v.ConsumeTokenIgnoreSeparator(lex.FuncKeywordToken)
	if err != nil {
		return err
	}

	node.Start = kw.View.Start
	return nil
}

func (p FunctionParser) parseBlockMaybe(v *TokenView, node *FunctionNode) {
	instructions, err := p.InstructionBlockParser.Parse(v)
	if err == nil {
		node.Instructions = &instructions
		node.End = node.Instructions.View().End
	} else {
		node.End = node.Declaration.View().End
	}
}

func (p FunctionParser) Parse(v *TokenView) (node FunctionNode, err ParsingError) {
	err = p.parseFunctionKeyword(v, &node)
	if err != nil {
		return
	}

	node.Declaration, err = p.FunctionDeclarationParser.Parse(v)
	if err != nil {
		return
	}

	p.parseBlockMaybe(v, &node)
	return
}
