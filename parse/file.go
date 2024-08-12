package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type FileNode struct {
	Functions []FunctionNode
	Types     []TypeDeclarationNode
}

func (n FileNode) countAllNodes() int {
	return len(n.Functions) + len(n.Types)
}

func (n FileNode) collectAllNodes() (nodes []Node) {
	nodes = make([]Node, 0, n.countAllNodes())
	for _, fun := range n.Functions {
		nodes = append(nodes, Node(fun))
	}
	for _, typ := range n.Types {
		nodes = append(nodes, Node(typ))
	}
	return
}

func (n FileNode) String(ctx source.SourceContext) (s string) {
	nodes := n.collectAllNodes()
	SortNodesBySourceOrder(nodes)
	// TODO: efficiency improvement, we do not need to sort here in O(n log n),
	// because each node type is already parsed in source order and stored in order.
	// we just need to merge the sorted lists in linear time.

	for i, node := range nodes {
		s += node.String(ctx)
		if i != len(nodes)-1 {
			s += "\n"
		}
	}

	return s
}

type FileParser struct {
	FunctionParser        FunctionParser
	TypeDeclarationParser TypeDeclarationParser
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
