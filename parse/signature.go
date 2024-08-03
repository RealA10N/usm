package parse

import (
	"usm/lex"
	"usm/source"
)

type SignatureNode struct {
	source.UnmanagedSourceView
	Identifier source.UnmanagedSourceView
	Parameters []ParameterNode
	Returns    []TypeNode
}

func (n SignatureNode) View() source.UnmanagedSourceView {
	return n.UnmanagedSourceView
}

func (n SignatureNode) String(ctx source.SourceContext) string {
	s := "def "
	for _, ret := range n.Returns {
		s += ret.String(ctx) + " "
	}

	s += string(n.Identifier.Raw(ctx))

	for _, arg := range n.Parameters {
		s += " " + arg.String(ctx)
	}

	return s
}

type SignatureParser struct {
	ParameterParser ParameterParser
	TypeParser      TypeParser
}

func (SignatureParser) parseDef(v *TokenView, node *SignatureNode) ParsingError {
	def, err := v.ConsumeToken(lex.DefToken)
	if err != nil {
		return err
	}

	node.Start = def.View.Start
	return nil
}

func (SignatureParser) parseIdentifier(v *TokenView, node *SignatureNode) ParsingError {
	id, err := v.ConsumeToken(lex.GlbToken)
	if err != nil {
		return err
	}

	node.Identifier = id.View
	return nil
}

func (SignatureParser) updateNodeViewEnd(node *SignatureNode) {
	if len(node.Parameters) > 0 {
		node.End = node.Parameters[len(node.Parameters)-1].View().End
	} else {
		node.End = node.Identifier.End
	}
}

func (p SignatureParser) Parse(v *TokenView) (node SignatureNode, err ParsingError) {
	err = p.parseDef(v, &node)
	if err != nil {
		return
	}

	node.Returns = ParseMany(p.TypeParser, v)

	err = p.parseIdentifier(v, &node)
	if err != nil {
		return
	}

	node.Parameters = ParseMany(p.ParameterParser, v)
	p.updateNodeViewEnd(&node)
	return
}
