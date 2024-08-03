package parse

type Parser[Node any] interface {
	Parse(v *TokenView) (Node, ParsingError)
}

func ParseMany[Node any](p Parser[Node], v *TokenView) (nodes []Node) {
	for {
		typ, err := p.Parse(v)
		if err != nil {
			return
		}
		nodes = append(nodes, typ)
	}
}
