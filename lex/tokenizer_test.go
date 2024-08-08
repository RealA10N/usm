package lex_test

import (
	"testing"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"

	"github.com/stretchr/testify/assert"
)

type tknDesc struct {
	txt string
	typ lex.TokenType
}

func assertExpectedTokens(t *testing.T, expected []tknDesc, actual []lex.Token, ctx source.SourceContext) {
	assert.Len(t, actual, len(expected))
	for i, act := range actual {
		exp := expected[i]
		actStr := string(act.View.Raw(ctx))
		assert.Equal(t, exp.txt, actStr, "expected '%s' got '%s'", exp.txt, actStr)
		assert.Equal(t, exp.typ, act.Type, "expected %s got %s", exp.typ, act.Type)
	}
}

func TestAddOne(t *testing.T) {
	code :=
		`function $32 @addOne $32 %x =
			%0 = add %x $32 #1
			ret %0
		`

	expected := []tknDesc{
		{"function", lex.FunctionKeywordToken},
		{"$32", lex.TypeToken},
		{"@addOne", lex.GlobalToken},
		{"$32", lex.TypeToken},
		{"%x", lex.RegisterToken},
		{"=", lex.EqualToken},
		{"", lex.SeparatorToken},
		{"%0", lex.RegisterToken},
		{"=", lex.EqualToken},
		{"add", lex.OperatorToken},
		{"%x", lex.RegisterToken},
		{"$32", lex.TypeToken},
		{"#1", lex.ImmediateToken},
		{"", lex.SeparatorToken},
		{"ret", lex.OperatorToken},
		{"%0", lex.RegisterToken},
		{"", lex.SeparatorToken},
	}

	view := source.NewSourceView(code)
	_, ctx := view.Detach()
	tkns, err := lex.NewTokenizer().Tokenize(view)

	assert.NoError(t, err)
	assertExpectedTokens(t, expected, tkns, ctx)
}

func TestPow(t *testing.T) {
	code :=
		`function $32 @pow $32 %base $32 %exp =
			jz %exp .end

		.recurse
			%base.new = mul %base %base
			%exp.new = shr %exp $32 #1
			%res.0 = call @pow %base.new %exp.new
			%exp.mod2 = and %exp $32 #1
			jz %exp.mod2 .even_base

		.odd_base
			%res.1 = mul %res.0 %base

		.even_base
			%res.2 = phi .odd_base %res.1 .recurse %res.0

		.end
			%res.3 = phi . %base .even_base %res.2
			ret %res.3
		`

	expected := []tknDesc{
		{"function", lex.FunctionKeywordToken},
		{"$32", lex.TypeToken},
		{"@pow", lex.GlobalToken},
		{"$32", lex.TypeToken},
		{"%base", lex.RegisterToken},
		{"$32", lex.TypeToken},
		{"%exp", lex.RegisterToken},
		{"=", lex.EqualToken},
		{"", lex.SeparatorToken},

		{"jz", lex.OperatorToken},
		{"%exp", lex.RegisterToken},
		{".end", lex.LabelToken},
		{"", lex.SeparatorToken},

		{".recurse", lex.LabelToken},
		{"", lex.SeparatorToken},

		{"%base.new", lex.RegisterToken},
		{"=", lex.EqualToken},
		{"mul", lex.OperatorToken},
		{"%base", lex.RegisterToken},
		{"%base", lex.RegisterToken},
		{"", lex.SeparatorToken},

		{"%exp.new", lex.RegisterToken},
		{"=", lex.EqualToken},
		{"shr", lex.OperatorToken},
		{"%exp", lex.RegisterToken},
		{"$32", lex.TypeToken},
		{"#1", lex.ImmediateToken},
		{"", lex.SeparatorToken},

		{"%res.0", lex.RegisterToken},
		{"=", lex.EqualToken},
		{"call", lex.OperatorToken},
		{"@pow", lex.GlobalToken},
		{"%base.new", lex.RegisterToken},
		{"%exp.new", lex.RegisterToken},
		{"", lex.SeparatorToken},

		{"%exp.mod2", lex.RegisterToken},
		{"=", lex.EqualToken},
		{"and", lex.OperatorToken},
		{"%exp", lex.RegisterToken},
		{"$32", lex.TypeToken},
		{"#1", lex.ImmediateToken},
		{"", lex.SeparatorToken},

		{"jz", lex.OperatorToken},
		{"%exp.mod2", lex.RegisterToken},
		{".even_base", lex.LabelToken},
		{"", lex.SeparatorToken},

		{".odd_base", lex.LabelToken},
		{"", lex.SeparatorToken},

		{"%res.1", lex.RegisterToken},
		{"=", lex.EqualToken},
		{"mul", lex.OperatorToken},
		{"%res.0", lex.RegisterToken},
		{"%base", lex.RegisterToken},
		{"", lex.SeparatorToken},

		{".even_base", lex.LabelToken},
		{"", lex.SeparatorToken},

		{"%res.2", lex.RegisterToken},
		{"=", lex.EqualToken},
		{"phi", lex.OperatorToken},
		{".odd_base", lex.LabelToken},
		{"%res.1", lex.RegisterToken},
		{".recurse", lex.LabelToken},
		{"%res.0", lex.RegisterToken},
		{"", lex.SeparatorToken},

		{".end", lex.LabelToken},
		{"", lex.SeparatorToken},

		{"%res.3", lex.RegisterToken},
		{"=", lex.EqualToken},
		{"phi", lex.OperatorToken},
		{".", lex.LabelToken},
		{"%base", lex.RegisterToken},
		{".even_base", lex.LabelToken},
		{"%res.2", lex.RegisterToken},
		{"", lex.SeparatorToken},

		{"ret", lex.OperatorToken},
		{"%res.3", lex.RegisterToken},
		{"", lex.SeparatorToken},
	}

	view := source.NewSourceView(code)
	_, ctx := view.Detach()
	tkns, err := lex.NewTokenizer().Tokenize(view)

	assert.NoError(t, err)
	assertExpectedTokens(t, expected, tkns, ctx)
}
