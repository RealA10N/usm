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
	Tokenize(core.SourceView) (TokenizeResult, error)
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

func (t tokenizer) Tokenize(view core.SourceView) (result TokenizeResult, err error) {
	for {
		addSep, comments := t.consumeWhitespace(&view)
		result.Comments = append(result.Comments, comments...)
		if addSep {
			result.Tokens = append(result.Tokens, Token{Type: SeparatorToken})
		}

		tkn, err := t.yieldToken(&view)
		if err != nil {
			break
		}

		result.Tokens = append(result.Tokens, tkn)
	}

	if view.Len() != 0 {
		return result, err
	}

	return result, nil
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

// consumeWhitespace consumes interleaved whitespace and ';' comments.
// Returns true if a newline was encountered, and any comments found.
func (tokenizer) consumeWhitespace(view *core.SourceView) (sawNewline bool, comments []Comment) {
	for {
		if consumeSpaces(view) {
			sawNewline = true
		}
		comment, ok := consumeComment(view)
		if !ok {
			return
		}
		comments = append(comments, comment)
	}
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
// Returns the comment and true, or the zero value and false if no comment was found.
func consumeComment(view *core.SourceView) (Comment, bool) {
	if !view.HasPrefix(core.NewSourceView(";")) {
		return Comment{}, false
	}
	idx := view.IndexFunc(func(r rune) bool { return r == '\n' })
	commentView, after := view.Partition(idx)
	detached, _ := commentView.Detach()
	*view = after
	return Comment{View: detached}, true
}

// Provided a boolean predicate, returns a new boolean predicate which yields
// opposite (not) values of the provided predicate.
func not[T any](f func(item T) bool) func(T) bool {
	return func(item T) bool {
		return !f(item)
	}
}
