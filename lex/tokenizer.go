package lex

import (
	"errors"
	"unicode"

	"alon.kr/x/usm/core"
)

type SpecificTokenizer interface {
	Tokenize(txt *core.SourceView) (Token, error)
}

type Tokenizer interface {
	Tokenize(core.SourceView) ([]Token, error)
}

type tokenizer struct {
	specificTokenizers []SpecificTokenizer
}

func NewTokenizer() Tokenizer {
	return tokenizer{
		specificTokenizers: []SpecificTokenizer{
			PrefixedTokenizer{"%", RegisterToken},
			PrefixedTokenizer{"$", TypeToken},
			PrefixedTokenizer{".", LabelToken},
			PrefixedTokenizer{"@", GlobalToken},
			PrefixedTokenizer{"#", ImmediateToken},
			PrefixedTokenizer{"*", PointerToken},
			PrefixedTokenizer{"^", RepeatToken},
			KeywordTokenizer{"=", EqualToken},
			KeywordTokenizer{"{", LeftCurlyBraceToken},
			KeywordTokenizer{"}", RightCurlyBraceToken},
			KeywordTokenizer{"func", FuncKeywordToken},
			KeywordTokenizer{"type", TypeKeywordToken},
			KeywordTokenizer{"var", VarKeywordToken},
			KeywordTokenizer{"const", ConstKeywordToken},
			WordTokenizer{OperatorToken},
		},
	}
}

func (t tokenizer) Tokenize(view core.SourceView) (tokens []Token, err error) {
	for {
		tokens = append(tokens, t.consumeWhitespace(&view)...)

		tkn, err := t.yieldToken(&view)
		if err != nil {
			break
		}

		tokens = append(tokens, tkn)
	}

	if view.Len() != 0 {
		return tokens, err
	}

	return tokens, nil
}

func (t tokenizer) yieldToken(view *core.SourceView) (tkn Token, err error) {
	for _, tokenParser := range t.specificTokenizers {
		tkn, err = tokenParser.Tokenize(view)
		if err == nil {
			return
		}
	}

	err = errors.New("unmatched subview")
	return
}

// consumeWhitespace consumes interleaved whitespace and ';' comments, returning
// the resulting token sequence. A SeparatorToken is emitted immediately after
// each run of whitespace that contains a newline; CommentTokens follow in source
// order. This ensures that a comment on a new line appears after the separator
// for that line, so the parser can distinguish inline comments from leading
// comments of the next node.
func (tokenizer) consumeWhitespace(view *core.SourceView) []Token {
	var tokens []Token
	for {
		if consumeSpaces(view) {
			tokens = append(tokens, Token{Type: SeparatorToken})
		}
		comment, ok := consumeComment(view)
		if !ok {
			break
		}
		tokens = append(tokens, comment)
	}
	return tokens
}

// consumeSpaces advances past leading whitespace and returns true if a newline was among them.
func consumeSpaces(view *core.SourceView) bool {
	idx := view.IndexFunc(not(unicode.IsSpace))
	before, after := view.Partition(idx)
	*view = after
	return before.Contains('\n')
}

// consumeComment advances past a ';'-style line comment if one is present.
// The trailing '\n' is left in the view so the caller can detect the line boundary.
// Returns a CommentToken and true, or the zero value and false if no comment was found.
func consumeComment(view *core.SourceView) (Token, bool) {
	if !view.HasPrefix(core.NewSourceView(";")) {
		return Token{}, false
	}
	idx := view.IndexFunc(func(r rune) bool { return r == '\n' })
	commentView, after := view.Partition(idx)
	detached, _ := commentView.Detach()
	*view = after
	return Token{Type: CommentToken, View: detached}, true
}

// Provided a boolean predicate, returns a new boolean predicate which yields
// opposite (not) values of the provided predicate.
func not[T any](f func(item T) bool) func(T) bool {
	return func(item T) bool {
		return !f(item)
	}
}
