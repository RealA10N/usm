package parse

import (
	"strings"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type InstructionNode struct {
	Operator  core.UnmanagedSourceView
	Arguments []ArgumentNode
	Targets   []ArgumentNode
	Labels    []LabelNode
	// LeadingComments holds whole-line comments before this instruction.
	LeadingComments []lex.Comment
	// TrailingComment holds the inline comment on the same line, after the last token.
	TrailingComment *lex.Comment
}

func (n *InstructionNode) attachLeadingComments(c []lex.Comment) {
	n.LeadingComments = c
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

func (n InstructionNode) stringTrailingComment(ctx *StringContext) string {
	if n.TrailingComment == nil {
		return ""
	}
	return " " + string(n.TrailingComment.View.Raw(ctx.SourceContext))
}

func (n InstructionNode) String(ctx *StringContext) string {
	return ctx.renderComments(n.LeadingComments) +
		n.stringLabels(ctx) +
		ctx.indent() + n.stringTargets(ctx) + string(n.Operator.Raw(ctx.SourceContext)) +
		n.stringArguments(ctx) + n.stringTrailingComment(ctx) + "\n"
}

type InstructionParser struct {
	LabelParser    Parser[LabelNode]
	ArgumentParser Parser[ArgumentNode]
}

func NewInstructionParser() InstructionParser {
	return InstructionParser{
		LabelParser:    NewLabelParser(),
		ArgumentParser: NewArgumentParser(),
	}
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
// > Lbl* (Arg+ Eql)? Opr? Arg* \n+
//
// Targets are parsed by greedily consuming arguments and committing them as
// targets only if an '=' token immediately follows. If '=' is absent the view
// is restored and the instruction is treated as having no targets.
//
// Note: leading whole-line comments are NOT consumed here; they are captured
// by BlockParser.parseBlockNodes and attached via attachLeadingComments.
func (p InstructionParser) Parse(v *TokenView) (node InstructionNode, err core.Result) {
	node.Labels, _ = ParseManyIgnoreSeparators(p.LabelParser, v)

	// Speculatively parse arguments as potential targets. Commit them as
	// targets only when followed by '='; otherwise restore the view.
	saved := *v
	candidates := ParseMany(p.ArgumentParser, v)
	if len(candidates) > 0 {
		if _, eqErr := v.ConsumeToken(lex.EqualToken); eqErr == nil {
			node.Targets = candidates
		} else {
			*v = saved
		}
	}

	err = p.parseOperator(v, &node)
	if err != nil {
		return
	}

	node.Arguments = ParseMany(p.ArgumentParser, v)
	node.TrailingComment = v.consumeTrailingComment()
	err = v.ConsumeAtLeastTokens(1, lex.SeparatorToken)
	return node, err
}
