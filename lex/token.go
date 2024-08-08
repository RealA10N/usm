package lex

import (
	"fmt"

	"alon.kr/x/usm/source"
)

type TokenType uint8

const (
	RegisterToken TokenType = iota
	TypeToken
	LabelToken
	GlobalToken
	ImmediateToken
	EqualToken
	FunctionKeywordToken
	GlobalKeywordToken
	TypeKeywordToken
	OperatorToken
	PointerToken
	RepeatToken
	SeparatorToken
)

var tokenNames = map[TokenType]string{
	RegisterToken:        "Register",
	TypeToken:            "Type",
	LabelToken:           "Label",
	GlobalToken:          "Global",
	ImmediateToken:       "Immediate",
	EqualToken:           "Equal",
	GlobalKeywordToken:   "Global Keyword",
	TypeKeywordToken:     "Type Keyword",
	FunctionKeywordToken: "Function Keyword",
	PointerToken:         "Pointer",
	RepeatToken:          "Repeat",
	OperatorToken:        "Operator",
	SeparatorToken:       "\n",
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
