package parse

import (
	"alon.kr/x/usm/lex"

	"alon.kr/x/view"
)

type TokenView struct{ view.View[lex.Token, uint32] }

func NewTokenView(tkns []lex.Token) TokenView {
	return TokenView{view.NewView[lex.Token, uint32](tkns)}
}

func (v *TokenView) ConsumeToken(expectedTypes ...lex.TokenType) (tkn lex.Token, perr ParsingError) {
	tknView, restView := v.Partition(1)
	tkn, err := tknView.At(0)

	if err != nil {
		perr = EofError{Expected: expectedTypes}
		return
	}

	gotExpected := false
	for _, expectedType := range expectedTypes {
		if tkn.Type == expectedType {
			gotExpected = true
			break
		}
	}

	if !gotExpected {
		perr = UnexpectedTokenError{Expected: expectedTypes, Actual: tkn}
		return
	}

	*v = TokenView{restView}
	return tkn, nil
}

// Consume as many tokens as possible greedly, until we receive an error.
// The error and token count are returned.
func (v *TokenView) ConsumeManyTokens(
	expectedTypes ...lex.TokenType,
) (count int, err ParsingError) {
	for ; ; count++ {
		_, err = v.ConsumeToken(expectedTypes...)
		if err != nil {
			return
		}
	}
}

// Consume as many tokens as possible greedly, until we recieve an error.
//
// If the number of tokens consumed is strictly less than the provided number,
// returns the underlying error. Otherwise, returns nil.
func (v *TokenView) ConsumeAtLeastTokens(
	atLeast int, expectedTypes ...lex.TokenType,
) (err ParsingError) {
	count, err := v.ConsumeManyTokens(expectedTypes...)
	if count < atLeast {
		return err
	} else {
		return nil
	}
}
