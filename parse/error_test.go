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

	expectedErr := `reached end of file`
	assert.Equal(t, expectedErr, err.Error(v.Ctx()))

	expectedHint := `expected <Equal>`
	assert.Equal(t, expectedHint, err.Hint(v.Ctx()))
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

	expectedErr := "reached end of file"
	assert.Equal(t, expectedErr, err.Error(v.Ctx()))

	expectedHint := "expected <Equal>, <Register>, <Pointer>"
	assert.Equal(t, expectedHint, err.Hint(v.Ctx()))
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

	expectedErr := `unexpected token <Register "%reg">`
	assert.Equal(t, expectedErr, err.Error(v.Ctx()))

	expectedHint := "expected <Equal>"
	assert.Equal(t, expectedHint, err.Hint(v.Ctx()))
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

	expectedErr := `unexpected token <Register "%reg">`
	assert.Equal(t, expectedErr, err.Error(v.Ctx()))

	expectedHint := "expected <Equal>, <Type>"
	assert.Equal(t, expectedHint, err.Hint(v.Ctx()))
}

func TestGenericUnexpectedError(t *testing.T) {
	v := core.NewSourceView("")
	err := parse.GenericUnexpectedError{"argument", v.Unmanaged()}

	expectedErr := "expected argument"
	assert.Equal(t, expectedErr, err.Error(v.Ctx()))

	assert.Empty(t, err.Hint(v.Ctx()))
}
