package parse

import (
	"usm/lex"
	"usm/source"
)

type InstructionNode struct {
	Operator  source.UnmanagedSourceView
	Arguments []CallerArgumentNode
	Targets   []RegisterNode
}

func (n InstructionNode) View() source.UnmanagedSourceView {
	first := n.Operator
	last := n.Operator

	if len(n.Targets) > 0 {
		first = n.Targets[0].View()
	}

	if len(n.Arguments) > 0 {
		last = n.Arguments[len(n.Arguments)-1].View()
	}

	return first.Merge(last)
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

type InstructionParser struct{}

func (InstructionParser) parseTargets(v *TokenView, node *InstructionNode) {
	for {
		reg, err := RegisterParser{}.Parse(v)
		if err != nil {
			return
		}
		node.Targets = append(node.Targets, reg)
	}
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

func (InstructionParser) parseArguments(v *TokenView, node *InstructionNode) ParsingError {
	for {
		arg, err := CallerArgumentParser{}.Parse(v)
		if err != nil {
			break
		}
		node.Arguments = append(node.Arguments, arg)
	}

	_, err := v.ConsumeToken(lex.SepToken)
	return err
}

func (p InstructionParser) Parse(v *TokenView) (node InstructionNode, err ParsingError) {
	p.parseTargets(v, &node)
	err = p.parseEquals(v, &node)
	if err != nil {
		return
	}

	err = p.parseOperator(v, &node)
	if err != nil {
		return
	}

	err = p.parseArguments(v, &node)
	return
}
