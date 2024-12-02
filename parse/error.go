package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

func stringManyTokenTypes(typs []lex.TokenType) (s string) {
	for i := 0; i < len(typs)-1; i++ {
		s += typs[i].String() + ", "
	}
	s += typs[len(typs)-1].String()
	return s
}

// MARK: UnexpectedTokenResult

func NewUnexpectedTokenResult(
	expected []lex.TokenType,
	actual lex.Token,
) core.Result {
	result := core.Result{{
		Type:     core.ErrorResult,
		Message:  "Unexpected token",
		Location: &actual.View,
	}}

	if len(expected) > 0 {
		result = append(result, core.ResultDetails{
			Type:    core.HintResult,
			Message: "Expected " + stringManyTokenTypes(expected),
		})
	}

	return result
}

// MARK: EofResult

func NewEofResult(expected []lex.TokenType) core.Result {
	v := core.NewEofUnmanagedSourceView()
	result := core.Result{{
		Type:     core.ErrorResult,
		Message:  "Reached end of file",
		Location: &v,
	}}

	if len(expected) > 0 {
		result = append(result, core.ResultDetails{
			Type:    core.HintResult,
			Message: "Expected " + stringManyTokenTypes(expected),
		})
	}

	return result
}
