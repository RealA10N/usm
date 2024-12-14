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

	if len(n.Labels) > 0 {
		v = v.MergeStart(n.Labels[0].View())
	}

	if len(n.Targets) > 0 {
		v = v.MergeStart(n.Targets[0].View())
	}

	if len(n.Arguments) > 0 {
		v = v.MergeEnd(n.Arguments[len(n.Arguments)-1].View())
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
	LabelParser    Parser[LabelNode]
	TargetParser   Parser[TargetNode]
	ArgumentParser Parser[ArgumentNode]
}

func NewInstructionParser() InstructionParser {
	return InstructionParser{
		LabelParser:    NewLabelParser(),
		TargetParser:   NewTargetParser(),
		ArgumentParser: NewArgumentParser(),
	}
}

func (InstructionParser) parseEquals(v *TokenView, node *InstructionNode) (err core.Result) {
	if len(node.Targets) > 0 {
		_, err = v.ConsumeToken(lex.EqualToken)
	}
	return
}

func (InstructionParser) parseOperator(v *TokenView, node *InstructionNode) core.Result {
	opr, err := v.ConsumeToken(lex.OperatorToken)
	if err != nil {
		// There is no operator: that is ok, the operator is the empty string

		nextToken, err := v.Front()
		if err != nil {
			return NewEofResult([]lex.TokenType{})
		}

		node.Operator = core.UnmanagedSourceView{
			Start: nextToken.View.Start,
			End:   nextToken.View.Start,
		}
		return nil
	}

	node.Operator = opr.View
	return nil
}

// Parsing of the following regular expression:
//
// > Lbl* (Reg+ Eql)? Opr? Arg+ !Arg \n+
func (p InstructionParser) Parse(v *TokenView) (node InstructionNode, err core.Result) {
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
	err = v.ConsumeAtLeastTokens(1, lex.SeparatorToken)
	return node, err
}
