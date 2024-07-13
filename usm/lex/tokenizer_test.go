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
	tokenizer := lex.Tokenizer{Reader: reader}
	gotTokens, err := tokenizer.Tokenize()

	assert.Nil(t, err)
	assert.Equal(t, gotTokens, expectedTokens)
}
