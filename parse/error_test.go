package parse_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func TestEofErrorOneExpected(t *testing.T) {
	v := core.NewSourceView("")
	err := parse.EofError{
		Expected: []lex.TokenType{lex.EqualToken},
	}

	expected := `reached end of file (expected <Equal>)`
	assert.Equal(t, expected, err.Error(v.Ctx()))
}

func TestEofErrorMultipleExpected(t *testing.T) {
	v := core.NewSourceView("")
	err := parse.EofError{
		Expected: []lex.TokenType{
			lex.EqualToken,
			lex.RegisterToken,
			lex.PointerToken,
		},
	}

	expected := `reached end of file (expected <Equal>, <Register>, <Pointer>)`
	assert.Equal(t, expected, err.Error(v.Ctx()))
}

func TestUnexpectedTokenOneExpected(t *testing.T) {
	v := core.NewSourceView("%reg")
	tkn := lex.Token{
		Type: lex.RegisterToken,
		View: v.Unmanaged(),
	}

	err := parse.UnexpectedTokenError{
		Expected: []lex.TokenType{lex.EqualToken},
		Actual:   tkn,
	}

	expected := `unexpected token <Register "%reg"> (expected <Equal>)`
	assert.Equal(t, expected, err.Error(v.Ctx()))
}

func TestUnexpectedTokenMultipleExpected(t *testing.T) {
	v := core.NewSourceView("%reg")
	tkn := lex.Token{
		Type: lex.RegisterToken,
		View: v.Unmanaged(),
	}

	err := parse.UnexpectedTokenError{
		Expected: []lex.TokenType{lex.EqualToken, lex.TypeToken},
		Actual:   tkn,
	}

	expected := `unexpected token <Register "%reg"> (expected <Equal>, <Type>)`
	assert.Equal(t, expected, err.Error(v.Ctx()))
}

func TestGenericUnexpectedError(t *testing.T) {
	v := core.NewSourceView("")
	err := parse.GenericUnexpectedError{"argument", v.Unmanaged()}
	expected := "expected argument"
	assert.Equal(t, expected, err.Error(v.Ctx()))
}
