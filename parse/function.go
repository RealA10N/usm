package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type FunctionNode struct {
	source.UnmanagedSourceView
	Declaration  FunctionDeclarationNode
	Instructions BlockNode[InstructionNode]
}

func (n FunctionNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n FunctionNode) String(ctx source.SourceContext) string {
	s := "function " + n.Declaration.String(ctx)
	if len(n.Instructions.Nodes) > 0 {
		s += " = " + n.Instructions.String(ctx)
	}
	return s
}

type FunctionParser struct {
	FunctionDeclarationParser FunctionDeclarationParser
	InstructionsParser        BlockParser[InstructionNode]
}

func (FunctionParser) String() string {
	return "function"
}

func (FunctionParser) parseFunctionKeyword(v *TokenView, node *FunctionNode) ParsingError {
	kw, err := v.ConsumeTokenIgnoreSeparator(lex.FunctionKeywordToken)
	if err != nil {
		return err
	}

	node.Start = kw.View.Start
	return nil
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

	_, err = v.ConsumeToken(lex.EqualToken)
	if err != nil {
		return
	}

	node.Instructions, err = p.InstructionsParser.Parse(v)
	node.End = node.Instructions.View().End
	return
}
