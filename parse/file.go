package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type FileNode struct {
	Functions []FunctionNode
	Types     []TypeDeclarationNode
}

func (n FileNode) View() (v source.UnmanagedSourceView) {
	if len(n.Functions) > 0 {
		v.Start = n.Functions[0].View().Start
		v.End = n.Functions[len(n.Functions)-1].View().End
	}
	return
}

func (n FileNode) String(ctx source.SourceContext) (s string) {
	l := len(n.Functions)
	if l == 0 {
		return
	}

	for i := 0; i < l-1; i++ {
		s += n.Functions[i].String(ctx) + "\n"
	}

	s += n.Functions[l-1].String(ctx)
	return s
}

type FileParser struct {
	FunctionParser        FunctionParser
	TypeDeclarationParser TypeDeclarationParser
}

func (FileParser) String() string {
	return "file"
}

func (p FileParser) parseNextNode(v *TokenView, node *FileNode) ParsingError {
	v.ConsumeManyTokens(lex.SeparatorToken)
	if v.Len() == 0 {
		return nil
	}

	tkn, err := v.PeekToken(lex.TopLevelTokens...)
	if err != nil {
		return err
	}

	switch tkn.Type {
	case lex.FuncKeywordToken:
		fun, err := p.FunctionParser.Parse(v)
		if err != nil {
			return err
		}
		node.Functions = append(node.Functions, fun)
	case lex.TypeKeywordToken:
		typ, err := p.TypeDeclarationParser.Parse(v)
		if err != nil {
			return err
		}
		node.Types = append(node.Types, typ)
	default:
		panic("unreachable")
	}

	return nil
}

func (p FileParser) Parse(v *TokenView) (node FileNode, err ParsingError) {
	for v.Len() > 0 {
		err = p.parseNextNode(v, &node)
		if err != nil {
			return
		}
	}
	return
}

func NewFileParser() FileParser {
	return FileParser{
		FunctionParser:        NewFunctionParser(),
		TypeDeclarationParser: NewTypeDeclarationParser(),
	}
}
