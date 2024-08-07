package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type InstructionNode struct {
	Operator  source.UnmanagedSourceView
	Arguments []ArgumentNode
	Targets   []RegisterNode
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
	op := string(n.Operator.Raw(ctx))
	return n.stringTargets(ctx) + op + n.stringArguments(ctx)
}

type InstructionParser struct {
	RegisterParser RegisterParser
	ArgumentParser ArgumentParser
}

func (InstructionParser) parseEquals(v *TokenView, node *InstructionNode) (err ParsingError) {
	if len(node.Targets) > 0 {
		_, err = v.ConsumeToken(lex.EqlToken)
	}
	return
}

func (InstructionParser) parseOperator(v *TokenView, node *InstructionNode) ParsingError {
	opr, err := v.ConsumeToken(lex.OprToken)
	node.Operator = opr.View
	return err
}

// Parsing of the following regular expression:
//
// > (Reg+ Eql)? Opr Arg+ !Arg
func (p InstructionParser) Parse(v *TokenView) (node InstructionNode, err ParsingError) {
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
	return
}
