package tokens

import (
	"usm/lex/base"
)

type LcrToken struct{}

func (LcrToken) String() string {
	return "<Lcr>"
}

type LcrTokenizer struct{}

func (LcrTokenizer) Tokenize(word string) (base.Token, error) {
	if word == "{" {
		return LcrToken{}, nil
	} else {
		return nil, base.ErrTokenNotMatched{Word: word}
	}
}
