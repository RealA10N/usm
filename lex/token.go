package lex

import (
	"fmt"
	"usm/source"
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

type Token struct {
	Type TokenType
	View source.UnmanagedSourceView
}

func (tkn Token) String(ctx source.SourceContext) string {
	typeName, ok := tokenNames[tkn.Type]
	if !ok {
		typeName = "?"
	}
	return fmt.Sprintf("<%s \"%s\">", typeName, string(tkn.View.Raw(ctx)))
}
