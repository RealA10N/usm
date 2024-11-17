package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

type ParsingError = core.UsmError

type UnexpectedTokenError struct {
	Expected []lex.TokenType
	Actual   lex.Token
}

func (e UnexpectedTokenError) Error(ctx core.SourceContext) string {
	return "unexpected token " + e.Actual.String(ctx)
}

func (e UnexpectedTokenError) Hint(core.SourceContext) string {
	if len(e.Expected) == 0 {
		return ""
	}

	return "expected " + stringManyTokenTypes(e.Expected)
}

func (e UnexpectedTokenError) Location() core.UnmanagedSourceView {
	return e.Actual.View
}

func (UnexpectedTokenError) IsInternalError() bool {
	return false
}

type EofError struct {
	Expected []lex.TokenType
}

func (e EofError) Error(core.SourceContext) string {
	return "reached end of file"
}

func (e EofError) Hint(core.SourceContext) string {
	if len(e.Expected) == 0 {
		return ""
	}

	return "expected " + stringManyTokenTypes(e.Expected)
}

func (e EofError) Location() core.UnmanagedSourceView {
	// TODO: This is a hack, but it's fine for now.
	return core.NewEofUnmanagedSourceView()
}

func (EofError) IsInternalError() bool {
	return false
}

func stringManyTokenTypes(typs []lex.TokenType) (s string) {
	for i := 0; i < len(typs)-1; i++ {
		s += typs[i].String() + ", "
	}
	s += typs[len(typs)-1].String()
	return s
}

type GenericUnexpectedError struct {
	Expected      string
	ErrorLocation core.UnmanagedSourceView
}

func (e GenericUnexpectedError) Error(core.SourceContext) string {
	return "expected " + e.Expected
}

func (e GenericUnexpectedError) Hint(core.SourceContext) string {
	return ""
}

func (e GenericUnexpectedError) Location() core.UnmanagedSourceView {
	return e.ErrorLocation
}

func (GenericUnexpectedError) IsInternalError() bool {
	return false
}
