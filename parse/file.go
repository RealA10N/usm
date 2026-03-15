package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type FileNode struct {
	Functions        []FunctionNode
	Types            []TypeDeclarationNode
	Constants        []ConstDeclarationNode
	Variables        []VarDeclarationNode
	TrailingComments []lex.Comment
}

func (n FileNode) View() core.UnmanagedSourceView {
	return core.NewFullUnmanagedSourceView()
}

func (n FileNode) countAllNodes() int {
	return len(n.Functions) + len(n.Types) + len(n.Constants) + len(n.Variables)
}

func (n FileNode) collectAllNodes() (nodes []Node) {
	nodes = make([]Node, 0, n.countAllNodes())
	for _, fun := range n.Functions {
		nodes = append(nodes, Node(fun))
	}
	for _, typ := range n.Types {
		nodes = append(nodes, Node(typ))
	}
	for _, constant := range n.Constants {
		nodes = append(nodes, Node(constant))
	}
	for _, variable := range n.Variables {
		nodes = append(nodes, Node(variable))
	}
	return
}

func (n FileNode) String(ctx *StringContext) (s string) {
	nodes := n.collectAllNodes()
	SortNodesBySourceOrder(nodes)
	// TODO: efficiency improvement, we do not need to sort here in O(n log n),
	// because each node type is already parsed in source order and stored in order.
	// we just need to merge the sorted lists in linear time.

	for i, node := range nodes {
		s += node.String(ctx) + "\n"
		if i != len(nodes)-1 {
			s += "\n"
		}
	}

	// Emit any trailing comments after the last node.
	if len(n.TrailingComments) > 0 {
		if len(nodes) > 0 {
			s += "\n"
		}
		s += ctx.renderComments(n.TrailingComments)
	}

	return s
}

type FileParser struct {
	FunctionParser         FunctionParser
	TypeDeclarationParser  TypeDeclarationParser
	ConstDeclarationParser Parser[ConstDeclarationNode]
	VarDeclarationParser   Parser[VarDeclarationNode]
}

func NewFileParser() FileParser {
	return FileParser{
		FunctionParser:         NewFunctionParser(),
		TypeDeclarationParser:  NewTypeDeclarationParser(),
		ConstDeclarationParser: NewConstDeclarationParser(),
		VarDeclarationParser:   NewVarDeclarationParser(),
	}
}

func (p FileParser) Parse(v *TokenView) (node FileNode, err core.Result) {
	for {
		pending := v.consumeLeadingComments()

		if v.Len() == 0 {
			node.TrailingComments = pending
			return
		}

		tkn, err := v.PeekToken(lex.TopLevelTokens...)
		if err != nil {
			return node, err
		}

		switch tkn.Type {
		case lex.FuncKeywordToken:
			fun, err := p.FunctionParser.Parse(v)
			if err != nil {
				return node, err
			}
			fun.LeadingComments = pending
			node.Functions = append(node.Functions, fun)
		case lex.TypeKeywordToken:
			typ, err := p.TypeDeclarationParser.Parse(v)
			if err != nil {
				return node, err
			}
			typ.LeadingComments = pending
			node.Types = append(node.Types, typ)
		case lex.ConstKeywordToken:
			constant, err := p.ConstDeclarationParser.Parse(v)
			if err != nil {
				return node, err
			}
			constant.LeadingComments = pending
			node.Constants = append(node.Constants, constant)
		case lex.VarKeywordToken:
			variable, err := p.VarDeclarationParser.Parse(v)
			if err != nil {
				return node, err
			}
			variable.LeadingComments = pending
			node.Variables = append(node.Variables, variable)
		default:
			return node, core.Result{
				{
					Type:     core.InternalErrorResult,
					Message:  "Unexpected token",
					Location: &tkn.View,
				},
			}
		}
	}
}
