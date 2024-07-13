package tokens

import (
	"usm/lex/base"
)

type DefToken struct {
}

func (DefToken) String() string {
	return "<Def>"
}

type DefTokenizer struct{}

func (DefTokenizer) Tokenize(word string) (base.Token, error) {
	if word == "def" {
		return DefToken{}, nil
	} else {
		return nil, base.ErrTokenNotMatched{Word: word}
	}
}
