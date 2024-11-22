package parse

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
)

// MARK: UnexpectedTokenError

type UnexpectedTokenError struct {
	Expected []lex.TokenType
	Actual   lex.Token
}

func (e UnexpectedTokenError) GetType() core.ResultType {
	return core.ErrorResult
}

func (e UnexpectedTokenError) GetMessage(ctx core.SourceContext) string {
	return "Unexpected token " + e.Actual.String(ctx)
}

func (e UnexpectedTokenError) GetLocation() *core.UnmanagedSourceView {
	return &e.Actual.View
}

func (e UnexpectedTokenError) GetNext() core.Result {
	if len(e.Expected) == 0 {
		return nil
	}

	return core.GenericResult{
		Type:    core.HintResult,
		Message: "Expected " + stringManyTokenTypes(e.Expected),
	}
}

// MARK: EofError

type EofError struct {
	Expected []lex.TokenType
}

func (e EofError) GetType() core.ResultType {
	return core.ErrorResult
}

func (e EofError) GetMessage(core.SourceContext) string {
	return "Reached end of file"
}

func (e EofError) GetLocation() *core.UnmanagedSourceView {
	// TODO: This is a hack, but it's fine for now.
	v := core.NewEofUnmanagedSourceView()
	return &v
}

func (e EofError) GetNext() core.Result {
	if len(e.Expected) == 0 {
		return nil
	}

	return core.GenericResult{
		Type:    core.HintResult,
		Message: "Expected " + stringManyTokenTypes(e.Expected),
	}
}

func stringManyTokenTypes(typs []lex.TokenType) (s string) {
	for i := 0; i < len(typs)-1; i++ {
		s += typs[i].String() + ", "
	}
	s += typs[len(typs)-1].String()
	return s
}
