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
	s := "unexpected token " + e.Actual.String(ctx)
	if len(e.Expected) > 0 {
		s += " (expected " + stringManyTokenTypes(e.Expected) + ")"
	}
	return s
}

func (e UnexpectedTokenError) Location() core.UnmanagedSourceView {
	return e.Actual.View
}

type EofError struct {
	Expected []lex.TokenType
}

func (e EofError) Error(core.SourceContext) string {
	s := "reached end of file"
	if len(e.Expected) > 0 {
		s += " (expected " + stringManyTokenTypes(e.Expected) + ")"
	}
	return s
}

func (e EofError) Location() core.UnmanagedSourceView {
	// TODO: This is a hack, but it's fine for now.
	return core.NewEofUnmanagedSourceView()
}

func stringManyTokenTypes(typs []lex.TokenType) (s string) {
	for i := 0; i < len(typs)-1; i++ {
		s += typs[i].String() + ", "
	}
	s += typs[len(typs)-1].String()
	return s
}

type GenericUnexpectedError struct {
	Expected       string
	SourceLocation core.UnmanagedSourceView
}

func (e GenericUnexpectedError) Error(core.SourceContext) string {
	return "expected " + e.Expected
}

func (e GenericUnexpectedError) Location() core.UnmanagedSourceView {
	return e.SourceLocation
}
