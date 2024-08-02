package parse

import (
	"usm/lex"
	"usm/source"

	"github.com/RealA10N/view"
)

type TokenView = view.View[lex.Token, uint32]
type UnmanagedTokenView = view.UnmanagedView[lex.Token, uint32]
type TokenViewContext = view.ViewContext[lex.Token]

type Node interface {
	// Return a reference to the node substring in the source code
	View() source.UnmanagedSourceView

	// Regenerate ("format") the code to a unique, single representation.
	String(ctx source.SourceContext) string
}

type NodeParser[T any] interface {
	Parse(view TokenView) (T, error)
}

// Parsing Utilities

func ConsumeToken(v *TokenView, expectedTypes ...lex.TokenType) (tkn lex.Token, perr ParsingError) {
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

	*v = restView
	return tkn, nil
}
