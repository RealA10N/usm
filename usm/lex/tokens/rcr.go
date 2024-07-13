package tokens

import (
	"usm/lex/base"
)

type RcrToken struct {
}

func (RcrToken) String() string {
	return "<Rcr>"
}

type RcrTokenizer struct {
}

func (RcrTokenizer) Tokenize(word string) (base.Token, error) {
	if word == "}" {
		return RcrToken{}, nil
	} else {
		return nil, base.ErrTokenNotMatched{Word: word}
	}
}
