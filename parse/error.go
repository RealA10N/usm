package parse

import (
	"fmt"
	"usm/lex"
	"usm/source"
)

type ParsingError interface {
	Error(source.SourceContext) string
}

type UnexpectedTokenError struct {
	Expected lex.TokenType
	Got      lex.Token
}

func (err UnexpectedTokenError) Error(ctx source.SourceContext) string {
	return fmt.Sprintf("expected %s token, but got %s", err.Expected.String(), err.Got.String(ctx))
}

type EofError struct {
	Expected lex.TokenType
}

func (err EofError) Error(source.SourceContext) string {
	return fmt.Sprintf("expected %s token, but file ended", err.Expected.String())
}
