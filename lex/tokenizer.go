package lex

import (
	"errors"
	"unicode"
)

type WordTokenizer interface {
	Tokenize(txt *SourceView) (Token, error)
}

type Tokenizer interface {
	Tokenize(SourceView) ([]Token, error)
}

type tokenizer struct {
	wordTokenizers []WordTokenizer
}

func NewTokenizer() Tokenizer {
	return tokenizer{
		wordTokenizers: []WordTokenizer{
			RegTokenizer{},
			TypTokenizer{},
			LblTokenizer{},
			GlbTokenizer{},
			ImmTokenizer{},
			LcrTokenizer{},
			RcrTokenizer{},
			EqlTokenizer{},
			OprTokenizer{},
		},
	}
}

func (t tokenizer) Tokenize(view SourceView) (tkns []Token, err error) {
	for {
		consumeWhitespace(&view)
		tkn, err := t.tokenizeWord(&view)
		if err != nil {
			break
		}
		tkns = append(tkns, tkn)
	}

	if view.Len() != 0 {
		return tkns, err
	}

	return tkns, nil
}

func (t tokenizer) tokenizeWord(view *SourceView) (tkn Token, err error) {
	for _, tokenParser := range t.wordTokenizers {
		tkn, err = tokenParser.Tokenize(view)
		if err == nil {
			return
		}
	}

	err = errors.New("unmatched subview")
	return
}

// Provided a boolean predicate, returns a new boolean predicate which yields
// opposite (not) values of the provided predicate.
func not[T any](f func(item T) bool) func(T) bool {
	return func(item T) bool {
		return !f(item)
	}
}

func consumeWhitespace(view *SourceView) {
	*view = view.Subview(view.Index(not(unicode.IsSpace)), view.Len())
}
