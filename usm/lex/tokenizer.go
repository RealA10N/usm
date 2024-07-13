package lex

import (
	"bufio"
	"io"
	"usm/lex/base"
	"usm/lex/tokens"
)

var wordTokenizers = []base.WordTokenizer{
	tokens.RegTokenizer{}, // %
	tokens.TypTokenizer{}, // $
	tokens.GlbTokenizer{}, // @
	tokens.LcrTokenizer{}, // {
	tokens.RcrTokenizer{}, // }
	tokens.EqlTokenizer{}, // =
	tokens.OprTokenizer{},
}

type Tokenizer struct {
	Reader io.Reader
}

func (tokenizer Tokenizer) Tokenize() ([]base.Token, error) {
	tokens := make([]base.Token, 0)

	scanner := bufio.NewScanner(tokenizer.Reader)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		token, err := tokenizer.tokenizeWord(word)
		if err != nil {
			// TODO: collect multiple errors and report all at once
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (Tokenizer) tokenizeWord(word string) (base.Token, error) {
	for _, tokenParser := range wordTokenizers {
		token, err := tokenParser.Tokenize(word)
		if err == nil {
			return token, nil
		}
	}

	return nil, base.ErrUnexpectedToken{Word: word}
}
