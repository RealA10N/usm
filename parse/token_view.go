package parse

import (
	"usm/lex"

	"github.com/RealA10N/view"
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

func (v *TokenView) ConsumeManyTokens(expectedTypes ...lex.TokenType) (tkns []lex.Token) {
	for {
		tkn, err := v.ConsumeToken(expectedTypes...)
		if err != nil {
			return
		}
		tkns = append(tkns, tkn)
	}
}
