package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type FunctionNode struct {
	source.UnmanagedSourceView
	Declaration  FunctionDeclarationNode
	Instructions []InstructionNode
}

func (n FunctionNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n FunctionNode) String(ctx source.SourceContext) string {
	s := "function " + n.Declaration.String(ctx)
	if len(n.Instructions) > 0 {
		s += " =\n"
		for _, inst := range n.Instructions {
			s += inst.String(ctx) + "\n"
		}
	}
	return s
}

type FunctionParser struct {
	FunctionDeclarationParser FunctionDeclarationParser
	InstructionParser         InstructionParser
}

func (FunctionParser) String() string {
	return "function"
}

func (FunctionParser) parseFunctionKeyword(v *TokenView, node *FunctionNode) ParsingError {
	kw, err := v.ConsumeToken(lex.FunctionKeywordToken)
	if err != nil {
		return err
	}

	node.Start = kw.View.Start
	return nil
}

func (p FunctionParser) parseInstructions(v *TokenView, node *FunctionNode) ParsingError {
	v.ConsumeManyTokens(lex.SeparatorToken)
	node.Instructions, _ = ParseManyConsumeSeparators(p.InstructionParser, v)

	if len(node.Instructions) > 0 {
		node.End = node.Instructions[len(node.Instructions)-1].View().End
	}

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

	err = p.parseInstructions(v, &node)
	return
}
