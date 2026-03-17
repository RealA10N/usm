package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

// MARK: Variable node (used as instruction argument)

type VariableNode struct{ TokenNode }
type VariableParser struct{ TokenParser[VariableNode] }

func VariableNodeCreator(tkn lex.Token) VariableNode {
	return VariableNode{TokenNode{tkn.View}}
}

func NewVariableParser() Parser[VariableNode] {
	return VariableParser{
		TokenParser: TokenParser[VariableNode]{
			Token:       lex.VariableToken,
			NodeCreator: VariableNodeCreator,
		},
	}
}

// MARK: Variable declaration node (appears in function body preamble)

type VariableDeclarationNode struct {
	Variable        VariableNode
	Type            TypeNode
	LeadingComments []lex.Comment
	TrailingComment *lex.Comment
}

func (n VariableDeclarationNode) View() core.UnmanagedSourceView {
	return n.Variable.View().MergeEnd(n.Type.View())
}

func (n *VariableDeclarationNode) attachLeadingComments(c []lex.Comment) {
	n.LeadingComments = c
}

func (n VariableDeclarationNode) String(ctx *StringContext) string {
	s := ctx.renderComments(n.LeadingComments)
	s += ctx.indent() + n.Variable.String(ctx) + " " + n.Type.String(ctx)
	if n.TrailingComment != nil {
		s += " " + string(n.TrailingComment.View.Raw(ctx.SourceContext))
	}
	return s + "\n"
}

// MARK: Variable declaration parser

type VariableDeclarationParser struct {
	VariableParser Parser[VariableNode]
	TypeParser     TypeParser
}

func NewVariableDeclarationParser() Parser[VariableDeclarationNode] {
	return VariableDeclarationParser{
		VariableParser: NewVariableParser(),
		TypeParser:     NewTypeParser(),
	}
}

func (p VariableDeclarationParser) Parse(v *TokenView) (node VariableDeclarationNode, err core.Result) {
	node.Variable, err = p.VariableParser.Parse(v)
	if err != nil {
		return
	}

	node.Type, err = p.TypeParser.Parse(v)
	if err != nil {
		return
	}

	node.TrailingComment = v.consumeTrailingComment()
	err = v.ConsumeAtLeastTokens(1, lex.SeparatorToken)
	return
}
