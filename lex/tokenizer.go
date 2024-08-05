package lex

import (
	"errors"
	"unicode"

	"alon.kr/x/usm/source"
)

type SpecificTokenizer interface {
	Tokenize(txt *source.SourceView) (Token, error)
}

type Tokenizer interface {
	Tokenize(source.SourceView) ([]Token, error)
}

type tokenizer struct {
	specificTokenizers []SpecificTokenizer
}

func NewTokenizer() Tokenizer {
	return tokenizer{
		specificTokenizers: []SpecificTokenizer{
			NewPrefixedTokenizer("%", RegToken),
			NewPrefixedTokenizer("$", TypToken),
			NewPrefixedTokenizer(".", LblToken),
			NewPrefixedTokenizer("@", GlbToken),
			NewPrefixedTokenizer("#", ImmToken),
			KeywordTokenizer{Keyword: "=", Token: EqlToken},
			KeywordTokenizer{Keyword: "{", Token: LcrToken},
			KeywordTokenizer{Keyword: "}", Token: RcrToken},
			KeywordTokenizer{Keyword: "def", Token: DefToken},
			WordTokenizer{Token: OprToken},
		},
	}
}

func (t tokenizer) Tokenize(view source.SourceView) (tkns []Token, err error) {
	for {
		addSep := t.consumeWhitespace(&view)
		if addSep {
			tkns = append(tkns, Token{Type: SepToken})
		}

		tkn, err := t.yieldToken(&view)
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

func (t tokenizer) yieldToken(view *source.SourceView) (tkn Token, err error) {
	for _, tokenParser := range t.specificTokenizers {
		tkn, err = tokenParser.Tokenize(view)
		if err == nil {
			return
		}
	}

	err = errors.New("unmatched subview")
	return
}

// Consume white spaces and return true if encounterd a newline.
func (tokenizer) consumeWhitespace(view *source.SourceView) bool {
	idx := view.Index(not(unicode.IsSpace))
	before, after := view.Partition(idx)
	*view = after
	return before.Contains('\n')
}

// Provided a boolean predicate, returns a new boolean predicate which yields
// opposite (not) values of the provided predicate.
func not[T any](f func(item T) bool) func(T) bool {
	return func(item T) bool {
		return !f(item)
	}
}
