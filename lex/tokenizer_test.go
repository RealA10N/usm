package lex_test

import (
	"strings"
	"testing"
	"usm/lex"
	"usm/lex/base"
	"usm/lex/tokens"

	"github.com/stretchr/testify/assert"
)

func TestAddOne(t *testing.T) {
	code :=
		`def $i32 @addOne $i32 %x {
			%0 = add %x #1
			ret %0
		}`

	expectedTokens := []base.Token{
		tokens.OprToken{Name: "def"},
		tokens.TypToken{Name: "i32"},
		tokens.GlbToken{Name: "addOne"},
		tokens.TypToken{Name: "i32"},
		tokens.RegToken{Name: "x"},
		tokens.LcrToken{},
		tokens.RegToken{Name: "0"},
		tokens.EqlToken{},
		tokens.OprToken{Name: "add"},
		tokens.RegToken{Name: "x"},
		tokens.ImmToken{Value: "1"},
		tokens.OprToken{Name: "ret"},
		tokens.RegToken{Name: "0"},
		tokens.RcrToken{},
	}

	reader := strings.NewReader(code)
	gotTokens, err := lex.Tokenizer{}.Tokenize(reader)

	assert.Nil(t, err)
	assert.Equal(t, expectedTokens, gotTokens)
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

	expectedTokens := []base.Token{
		tokens.OprToken{Name: "def"},
		tokens.TypToken{Name: "u32"},
		tokens.GlbToken{Name: "pow"},
		tokens.TypToken{Name: "u32"},
		tokens.RegToken{Name: "base"},
		tokens.TypToken{Name: "u32"},
		tokens.RegToken{Name: "exp"},
		tokens.LcrToken{},

		tokens.OprToken{Name: "jz"},
		tokens.RegToken{Name: "exp"},
		tokens.LblToken{Name: "end"},

		tokens.LblToken{Name: "recurse"},

		tokens.RegToken{Name: "base.new"},
		tokens.EqlToken{},
		tokens.OprToken{Name: "mul"},
		tokens.RegToken{Name: "base"},
		tokens.RegToken{Name: "base"},

		tokens.RegToken{Name: "exp.new"},
		tokens.EqlToken{},
		tokens.OprToken{Name: "shr"},
		tokens.RegToken{Name: "exp"},
		tokens.ImmToken{Value: "1"},

		tokens.RegToken{Name: "res.0"},
		tokens.EqlToken{},
		tokens.OprToken{Name: "call"},
		tokens.GlbToken{Name: "pow"},
		tokens.RegToken{Name: "base.new"},
		tokens.RegToken{Name: "exp.new"},

		tokens.RegToken{Name: "exp.mod2"},
		tokens.EqlToken{},
		tokens.OprToken{Name: "and"},
		tokens.RegToken{Name: "exp"},
		tokens.ImmToken{Value: "1"},

		tokens.OprToken{Name: "jz"},
		tokens.RegToken{Name: "exp.mod2"},
		tokens.LblToken{Name: "even_base"},

		tokens.LblToken{Name: "odd_base"},

		tokens.RegToken{Name: "res.1"},
		tokens.EqlToken{},
		tokens.OprToken{Name: "mul"},
		tokens.RegToken{Name: "res.0"},
		tokens.RegToken{Name: "base"},

		tokens.LblToken{Name: "even_base"},

		tokens.RegToken{Name: "res.2"},
		tokens.EqlToken{},
		tokens.OprToken{Name: "phi"},
		tokens.LblToken{Name: "odd_base"},
		tokens.RegToken{Name: "res.1"},
		tokens.LblToken{Name: "recurse"},
		tokens.RegToken{Name: "res.0"},

		tokens.LblToken{Name: "end"},

		tokens.RegToken{Name: "res.3"},
		tokens.EqlToken{},
		tokens.OprToken{Name: "phi"},
		tokens.LblToken{Name: ""},
		tokens.RegToken{Name: "base"},
		tokens.LblToken{Name: "even_base"},
		tokens.RegToken{Name: "res.2"},

		tokens.OprToken{Name: "ret"},
		tokens.RegToken{Name: "res.3"},

		tokens.RcrToken{},
	}

	reader := strings.NewReader(code)
	gotTokens, err := lex.Tokenizer{}.Tokenize(reader)

	assert.Nil(t, err)
	assert.Equal(t, expectedTokens, gotTokens)
}
