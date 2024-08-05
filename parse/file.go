package parse

import "alon.kr/x/usm/source"

type FileNode struct {
	Functions []FunctionNode
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
	FunctionParser FunctionParser
}

func (p FileParser) Parse(v *TokenView) (node FileNode, err ParsingError) {
	node.Functions, _ = ParseManyConsumeSeparators(p.FunctionParser, v)
	return
}
