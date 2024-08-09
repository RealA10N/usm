package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

// TODO add function labels before each instruction.

type InstructionNode struct {
	Operator  source.UnmanagedSourceView
	Arguments []ArgumentNode
	Targets   []RegisterNode
	Labels    []LabelNode
}

func (n InstructionNode) View() (v source.UnmanagedSourceView) {
	v = n.Operator

	if len(n.Targets) > 0 {
		v.Start = n.Targets[0].View().Start
	}

	if len(n.Arguments) > 0 {
		v.End = n.Arguments[len(n.Arguments)-1].View().End
	}

	return
}

func (n InstructionNode) stringLabels(ctx source.SourceContext) (s string) {
	if len(n.Labels) == 0 {
		return
	}

	for _, lbl := range n.Labels {
		s += lbl.String(ctx) + " "
	}

	return
}

func (n InstructionNode) stringArguments(ctx source.SourceContext) (s string) {
	if len(n.Arguments) == 0 {
		return
	}

	for _, arg := range n.Arguments {
		s += " " + arg.String(ctx)
	}

	return
}

func (n InstructionNode) stringTargets(ctx source.SourceContext) (s string) {
	if len(n.Targets) == 0 {
		return
	}

	for _, tgt := range n.Targets {
		s += tgt.String(ctx) + " "
	}

	s += "= "
	return
}

func (n InstructionNode) String(ctx source.SourceContext) string {
	labels := n.stringLabels(ctx)
	targets := n.stringTargets(ctx)
	op := string(n.Operator.Raw(ctx))
	arguments := n.stringArguments(ctx)
	return labels + targets + op + arguments
}

type InstructionParser struct {
	LabelParser    LabelParser
	RegisterParser RegisterParser
	ArgumentParser ArgumentParser
}

func (InstructionParser) parseEquals(v *TokenView, node *InstructionNode) (err ParsingError) {
	if len(node.Targets) > 0 {
		_, err = v.ConsumeToken(lex.EqualToken)
	}
	return
}

func (InstructionParser) parseOperator(v *TokenView, node *InstructionNode) ParsingError {
	opr, err := v.ConsumeToken(lex.OperatorToken)
	node.Operator = opr.View
	return err
}

// Parsing of the following regular expression:
//
// > Lbl* (Reg+ Eql)? Opr Arg+ !Arg
func (p InstructionParser) Parse(v *TokenView) (node InstructionNode, err ParsingError) {
	node.Labels, _ = ParseManyIgnoreSeparators(p.LabelParser, v)
	node.Targets = ParseMany(p.RegisterParser, v)

	err = p.parseEquals(v, &node)
	if err != nil {
		return
	}

	err = p.parseOperator(v, &node)
	if err != nil {
		return
	}

	node.Arguments = ParseMany(p.ArgumentParser, v)
	return node, nil
}
