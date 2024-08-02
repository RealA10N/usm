package parse

import (
	"usm/lex"
	"usm/source"

	"github.com/RealA10N/view"
)

type TokenView = view.View[lex.Token, uint32]
type UnmanagedTokenView = view.UnmanagedView[lex.Token, uint32]
type TokenViewContext = view.ViewContext[lex.Token]

type Node struct {
	View source.UnmanagedSourceView
}

type NodeParser[T any] interface {
	Parse(view TokenView) (T, error)
}
