package tokens

import (
	"usm/lex/base"
)

type EqlToken struct {
}

func (EqlToken) String() string {
	return "<Eql>"
}

type EqlTokenizer struct {
}

func (EqlTokenizer) Tokenize(word string) (base.Token, error) {
	if word == "=" {
		return RcrToken{}, nil
	} else {
		return nil, base.ErrTokenNotMatched{Word: word}
	}
}
