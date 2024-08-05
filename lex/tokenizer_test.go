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
		`def $i32 @addOne $i32 %x {
			%0 = add %x #1
			ret %0
		}`

	expected := []tknDesc{
		tknDesc{"def", lex.DefToken},
		tknDesc{"$i32", lex.TypToken},
		tknDesc{"@addOne", lex.GlbToken},
		tknDesc{"$i32", lex.TypToken},
		tknDesc{"%x", lex.RegToken},
		tknDesc{"{", lex.LcrToken},
		tknDesc{"", lex.SepToken},
		tknDesc{"%0", lex.RegToken},
		tknDesc{"=", lex.EqlToken},
		tknDesc{"add", lex.OprToken},
		tknDesc{"%x", lex.RegToken},
		tknDesc{"#1", lex.ImmToken},
		tknDesc{"", lex.SepToken},
		tknDesc{"ret", lex.OprToken},
		tknDesc{"%0", lex.RegToken},
		tknDesc{"", lex.SepToken},
		tknDesc{"}", lex.RcrToken},
	}

	view := source.NewSourceView(code)
	_, ctx := view.Detach()
	tkns, err := lex.NewTokenizer().Tokenize(view)

	assert.NoError(t, err)
	assertExpectedTokens(t, expected, tkns, ctx)
}

func TestPow(t *testing.T) {
	code :=
		`def $u32 @pow $u32 %base $u32 %exp {
			jz %exp .end

		.recurse
			%base.new = mul %base %base
			%exp.new = shr %exp #1
			%res.0 = call @pow %base.new %exp.new
			%exp.mod2 = and %exp #1
			jz %exp.mod2 .even_base

		.odd_base
			%res.1 = mul %res.0 %base

		.even_base
			%res.2 = phi .odd_base %res.1 .recurse %res.0

		.end
			%res.3 = phi . %base .even_base %res.2
			ret %res.3
		}`

	expected := []tknDesc{
		tknDesc{"def", lex.DefToken},
		tknDesc{"$u32", lex.TypToken},
		tknDesc{"@pow", lex.GlbToken},
		tknDesc{"$u32", lex.TypToken},
		tknDesc{"%base", lex.RegToken},
		tknDesc{"$u32", lex.TypToken},
		tknDesc{"%exp", lex.RegToken},
		tknDesc{"{", lex.LcrToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"jz", lex.OprToken},
		tknDesc{"%exp", lex.RegToken},
		tknDesc{".end", lex.LblToken},
		tknDesc{"", lex.SepToken},

		tknDesc{".recurse", lex.LblToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"%base.new", lex.RegToken},
		tknDesc{"=", lex.EqlToken},
		tknDesc{"mul", lex.OprToken},
		tknDesc{"%base", lex.RegToken},
		tknDesc{"%base", lex.RegToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"%exp.new", lex.RegToken},
		tknDesc{"=", lex.EqlToken},
		tknDesc{"shr", lex.OprToken},
		tknDesc{"%exp", lex.RegToken},
		tknDesc{"#1", lex.ImmToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"%res.0", lex.RegToken},
		tknDesc{"=", lex.EqlToken},
		tknDesc{"call", lex.OprToken},
		tknDesc{"@pow", lex.GlbToken},
		tknDesc{"%base.new", lex.RegToken},
		tknDesc{"%exp.new", lex.RegToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"%exp.mod2", lex.RegToken},
		tknDesc{"=", lex.EqlToken},
		tknDesc{"and", lex.OprToken},
		tknDesc{"%exp", lex.RegToken},
		tknDesc{"#1", lex.ImmToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"jz", lex.OprToken},
		tknDesc{"%exp.mod2", lex.RegToken},
		tknDesc{".even_base", lex.LblToken},
		tknDesc{"", lex.SepToken},

		tknDesc{".odd_base", lex.LblToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"%res.1", lex.RegToken},
		tknDesc{"=", lex.EqlToken},
		tknDesc{"mul", lex.OprToken},
		tknDesc{"%res.0", lex.RegToken},
		tknDesc{"%base", lex.RegToken},
		tknDesc{"", lex.SepToken},

		tknDesc{".even_base", lex.LblToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"%res.2", lex.RegToken},
		tknDesc{"=", lex.EqlToken},
		tknDesc{"phi", lex.OprToken},
		tknDesc{".odd_base", lex.LblToken},
		tknDesc{"%res.1", lex.RegToken},
		tknDesc{".recurse", lex.LblToken},
		tknDesc{"%res.0", lex.RegToken},
		tknDesc{"", lex.SepToken},

		tknDesc{".end", lex.LblToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"%res.3", lex.RegToken},
		tknDesc{"=", lex.EqlToken},
		tknDesc{"phi", lex.OprToken},
		tknDesc{".", lex.LblToken},
		tknDesc{"%base", lex.RegToken},
		tknDesc{".even_base", lex.LblToken},
		tknDesc{"%res.2", lex.RegToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"ret", lex.OprToken},
		tknDesc{"%res.3", lex.RegToken},
		tknDesc{"", lex.SepToken},

		tknDesc{"}", lex.RcrToken},
	}

	view := source.NewSourceView(code)
	_, ctx := view.Detach()
	tkns, err := lex.NewTokenizer().Tokenize(view)

	assert.NoError(t, err)
	assertExpectedTokens(t, expected, tkns, ctx)
}
