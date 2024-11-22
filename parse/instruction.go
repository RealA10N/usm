package parse

import (
	"strings"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type InstructionNode struct {
	Operator  core.UnmanagedSourceView
	Arguments []ArgumentNode
	Targets   []TargetNode
	Labels    []LabelNode
}

func (n InstructionNode) View() (v core.UnmanagedSourceView) {
	v = n.Operator

	if len(n.Targets) > 0 {
		v.Start = n.Targets[0].View().Start
	}

	if len(n.Arguments) > 0 {
		v.End = n.Arguments[len(n.Arguments)-1].View().End
	}

	return
}

func (n InstructionNode) stringLabels(ctx *StringContext) (s string) {
	prefix := strings.Repeat("\t", max(0, ctx.Indent-1))
	for _, label := range n.Labels {
		s += prefix + label.String(ctx) + "\n"
	}
	return
}

func (n InstructionNode) stringArguments(ctx *StringContext) (s string) {
	for _, arg := range n.Arguments {
		s += " " + arg.String(ctx)
	}

	return
}

func (n InstructionNode) stringTargets(ctx *StringContext) (s string) {
	if len(n.Targets) == 0 {
		return
	}

	for _, tgt := range n.Targets {
		s += tgt.String(ctx) + " "
	}

	s += "= "
	return
}

func (n InstructionNode) String(ctx *StringContext) string {
	labels := n.stringLabels(ctx)
	prefix := strings.Repeat("\t", ctx.Indent)
	targets := n.stringTargets(ctx)
	op := string(n.Operator.Raw(ctx.SourceContext))
	arguments := n.stringArguments(ctx)
	return labels + prefix + targets + op + arguments + "\n"
}

type InstructionParser struct {
	LabelParser
	TargetParser
	ArgumentParser
}

func NewInstructionParser() InstructionParser {
	return InstructionParser{
		LabelParser:    NewLabelParser(),
		TargetParser:   NewTargetParser(),
		ArgumentParser: NewArgumentParser(),
	}
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
	node.Targets = ParseMany(p.TargetParser, v)

	err = p.parseEquals(v, &node)
	if err != nil {
		return
	}

	err = p.parseOperator(v, &node)
	if err != nil {
		return
	}

	node.Arguments = ParseMany(p.ArgumentParser, v)
	v.ConsumeManyTokens(lex.SeparatorToken)
	return node, nil
}
