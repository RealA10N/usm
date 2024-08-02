package parse

import (
	"usm/lex"
	"usm/source"

	"github.com/RealA10N/view"
)

type TokenView = view.View[lex.Token, uint32]
type UnmanagedTokenView = view.UnmanagedView[lex.Token, uint32]
type TokenViewContext = view.ViewContext[lex.Token]

func SourceViewFromBoundaryTokens(first, last lex.Token) source.UnmanagedSourceView {
	return source.UnmanagedSourceView{
		Start: first.View.Start,
		End:   last.View.End,
	}
}

type NodeParser[T any] interface {
	Parse(view TokenView) (T, error)
}

// Parsing Utilities

func ConsumeToken(v *TokenView, typ lex.TokenType) (tkn lex.Token, perr ParsingError) {
	tknView, restView := v.Partition(1)
	tkn, err := tknView.At(0)

	if err != nil {
		perr = EofError{Expected: typ}
		return
	}

	if tkn.Type != typ {
		perr = UnexpectedTokenError{Expected: typ, Got: tkn}
		return
	}

	*v = restView
	return tkn, nil
}
