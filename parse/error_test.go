package parse_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

// TODO: add GetLocation tests.

func TestEofErrorOneExpected(t *testing.T) {
	v := core.NewSourceView("")
	err := core.Result(
		parse.EofError{
			Expected: []lex.TokenType{lex.EqualToken},
		},
	)

	assert.Equal(t, "Reached end of file", err.GetMessage(v.Ctx()))
	assert.EqualValues(t, core.ErrorResult, err.GetType())

	hint := err.GetNext()
	assert.NotNil(t, hint)

	assert.EqualValues(t, core.HintResult, hint.GetType())
	assert.Equal(t, "Expected <Equal>", hint.GetMessage(v.Ctx()))
	assert.Nil(t, hint.GetNext())
}

func TestEofErrorMultipleExpected(t *testing.T) {
	v := core.NewSourceView("")
	err := core.Result(
		parse.EofError{
			Expected: []lex.TokenType{
				lex.EqualToken,
				lex.RegisterToken,
				lex.PointerToken,
			},
		},
	)

	assert.Equal(t, "Reached end of file", err.GetMessage(v.Ctx()))
	assert.EqualValues(t, core.ErrorResult, err.GetType())

	hint := err.GetNext()
	assert.NotNil(t, hint)

	assert.EqualValues(t, core.HintResult, hint.GetType())
	assert.Equal(t, "Expected <Equal>, <Register>, <Pointer>", hint.GetMessage(v.Ctx()))
	assert.Nil(t, hint.GetNext())
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

	assert.Equal(t, `Unexpected token <Register "%reg">`, err.GetMessage(v.Ctx()))
	assert.EqualValues(t, core.ErrorResult, err.GetType())

	hint := err.GetNext()
	assert.NotNil(t, hint)

	assert.EqualValues(t, core.HintResult, hint.GetType())
	assert.Equal(t, "Expected <Equal>", hint.GetMessage(v.Ctx()))
	assert.Nil(t, hint.GetNext())
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
	assert.Equal(t, `Unexpected token <Register "%reg">`, err.GetMessage(v.Ctx()))
	assert.EqualValues(t, core.ErrorResult, err.GetType())

	hint := err.GetNext()
	assert.NotNil(t, hint)

	assert.EqualValues(t, core.HintResult, hint.GetType())
	assert.Equal(t, "Expected <Equal>, <Type>", hint.GetMessage(v.Ctx()))
	assert.Nil(t, hint.GetNext())
}
