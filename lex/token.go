package lex

import (
	"fmt"

	"github.com/RealA10N/view"
)

type TokenType uint8

const (
	RegToken TokenType = iota
	TypToken
	LblToken
	GlbToken
	ImmToken
	LcrToken
	RcrToken
	EqlToken
	OprToken
)

var tokenNames = map[TokenType]string{
	RegToken: "Register",
	TypToken: "Type",
	LblToken: "Label",
	GlbToken: "Global",
	ImmToken: "Immediate",
	LcrToken: "Left Curly Brace",
	RcrToken: "Right Curly Brace",
	EqlToken: "Equal",
	OprToken: "Operator",
}

func (tkn TokenType) String() string {
	name, ok := tokenNames[tkn]
	if !ok {
		name = "?"
	}
	return fmt.Sprintf("<%s>", name)
}

type SourceView = view.View[rune, uint32]
type UnmanagedSourceView = view.UnmanagedView[rune, uint32]
type SourceContext = view.ViewContext[rune]

func NewSourceView(data string) SourceView {
	return view.NewView[rune, uint32]([]rune(data))
}

type Token struct {
	Type TokenType
	View UnmanagedSourceView
}

func (tkn Token) String(ctx SourceContext) string {
	typeName, ok := tokenNames[tkn.Type]
	if !ok {
		typeName = "?"
	}
	return fmt.Sprintf("<%s \"%s\">", typeName, string(tkn.View.Raw(ctx)))
}
