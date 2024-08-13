package parse

type ConstNode = GlobalDeclarationNode

type ConstParser struct {
	GlobalDeclarationParser GlobalDeclarationParser
}

func (p ConstParser) Parse(v *TokenView) (
	node ConstNode,
	err ParsingError,
) {
	node, err = p.GlobalDeclarationParser.Parse(v)
	return
}
