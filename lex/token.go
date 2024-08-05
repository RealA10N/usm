package lex

import (
	"fmt"

	"alon.kr/x/usm/source"
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
	DefToken
	OprToken
	SepToken
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
	DefToken: "Define",
	OprToken: "Operator",
	SepToken: `\n`,
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

	if tkn.View.Len() > 0 {
		return fmt.Sprintf(`<%s "%s">`, typeName, string(tkn.View.Raw(ctx)))
	} else {
		return fmt.Sprintf(typeName)
	}
}
