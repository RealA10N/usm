package lex

import (
	"fmt"

	"alon.kr/x/usm/core"
)

type TokenType uint8

const (
	RegisterToken TokenType = iota
	TypeToken
	LabelToken
	GlobalToken
	ImmediateToken
	EqualToken
	LeftCurlyBraceToken
	RightCurlyBraceToken
	FuncKeywordToken
	TypeKeywordToken
	VarKeywordToken
	ConstKeywordToken
	PointerToken
	RepeatToken
	OperatorToken
	SeparatorToken
)

var TopLevelTokens = []TokenType{
	SeparatorToken,
	FuncKeywordToken,
	TypeKeywordToken,
	VarKeywordToken,
	ConstKeywordToken,
}

var tokenNames = map[TokenType]string{
	RegisterToken:        "Register",
	TypeToken:            "Type",
	LabelToken:           "Label",
	GlobalToken:          "Global",
	ImmediateToken:       "Immediate",
	EqualToken:           "Equal",
	LeftCurlyBraceToken:  "Left Curly Brace",
	RightCurlyBraceToken: "Right Curly Brace",
	FuncKeywordToken:     "Func Keyword",
	TypeKeywordToken:     "Type Keyword",
	VarKeywordToken:      "Var Keyword",
	ConstKeywordToken:    "Const Keyword",
	PointerToken:         "Pointer",
	RepeatToken:          "Repeat",
	OperatorToken:        "Operator",
	SeparatorToken:       `\n`,
}

func (tkn TokenType) String() string {
	name, ok := tokenNames[tkn]
	if !ok {
		return "?"
	}
	return name
}

type Token struct {
	Type TokenType
	View core.UnmanagedSourceView
}

func (tkn Token) String(ctx core.SourceContext) string {
	typeName := tkn.Type.String()

	if tkn.View.Len() > 0 {
		return fmt.Sprintf(`<%s "%s">`, typeName, string(tkn.View.Raw(ctx)))
	} else {
		return fmt.Sprintf(`<%s>`, typeName)
	}
}
