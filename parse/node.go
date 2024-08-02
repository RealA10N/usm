package parse

import (
	"usm/lex"
	"usm/source"

	"github.com/RealA10N/view"
)

type Node interface {
	// Return a reference to the node substring in the source code
	View() source.UnmanagedSourceView

	// Regenerate ("format") the code to a unique, single representation.
	String(ctx source.SourceContext) string
}

type TokenView struct{ view.View[lex.Token, uint32] }

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
