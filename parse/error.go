package parse

import (
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/source"
)

type ParsingError interface {
	Error(source.SourceContext) string
}

type UnexpectedTokenError struct {
	Expected []lex.TokenType
	Actual   lex.Token
}

func (e UnexpectedTokenError) Error(ctx source.SourceContext) string {
	s := "got token " + e.Actual.String(ctx)
	if len(e.Expected) > 0 {
		s += " (expected " + stringManyTokenTypes(e.Expected) + ")"
	}
	return s
}

type EofError struct {
	Expected []lex.TokenType
}

func (e EofError) Error(source.SourceContext) string {
	s := "reached end of file"
	if len(e.Expected) > 0 {
		s += " (expected " + stringManyTokenTypes(e.Expected) + ")"
	}
	return s
}

func stringManyTokenTypes(typs []lex.TokenType) (s string) {
	for i := 0; i < len(typs)-1; i++ {
		s += typs[i].String() + ", "
	}
	s += typs[len(typs)-1].String()
	return s
}

type GenericUnexpectedError struct {
	Expected string
}

func (e GenericUnexpectedError) Error(ctx source.SourceContext) string {
	return "expected " + e.Expected
}
