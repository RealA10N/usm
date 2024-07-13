package parser

import (
	"bufio"
	"io"
	"strings"
	"usm/tokens"
)

type Parser struct {
	Reader io.Reader
}

func (parser Parser) Parse() ([]tokens.Token, error) {
	tokens := make([]tokens.Token, 0)

	scanner := bufio.NewScanner(parser.Reader)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		token, err := parser.parseToken(word)
		if err != nil {
			// TODO: collect multiple errors and report all at once
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

var tokenizers = []tokens.Tokenizer{
	tokens.DefTokenizer{},
	tokens.RegisterTokenizer{},
	tokens.TypeTokenizer{},
	tokens.GlobalTokenizer{},
	tokens.LcurlTokenizer{},
	tokens.RcurlTokenizer{},
}

func (parser Parser) parseToken(word string) (tokens.Token, error) {
	for _, tokenParser := range tokenizers {
		token, err := tokenParser.Tokenize(word)
		if err == nil {
			return token, nil
		}
	}

	return nil, tokens.ErrUnexpectedToken{Word: word}
}
