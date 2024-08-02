package lex

import (
	"errors"
	"unicode"
	"usm/source"
)

type WordTokenizer interface {
	Tokenize(txt *source.SourceView) (Token, error)
}

type Tokenizer interface {
	Tokenize(source.SourceView) ([]Token, error)
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
			KeywordTokenizer{Keyword: "=", Token: EqlToken},
			KeywordTokenizer{Keyword: "{", Token: LcrToken},
			KeywordTokenizer{Keyword: "}", Token: RcrToken},
			KeywordTokenizer{Keyword: "def", Token: DefToken},
			OprTokenizer{},
		},
	}
}

func (t tokenizer) Tokenize(view source.SourceView) (tkns []Token, err error) {
	for {
		addSep := consumeWhitespace(&view)
		if addSep {
			tkns = append(tkns, Token{Type: SepToken})
		}
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

func (t tokenizer) tokenizeWord(view *source.SourceView) (tkn Token, err error) {
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

// consume white spaces and return true if encounterd a newline.
func consumeWhitespace(view *source.SourceView) bool {
	idx := view.Index(not(unicode.IsSpace))
	before, after := view.Partition(idx)
	*view = after
	return before.Contains('\n')
}
