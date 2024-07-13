package tokens

type DefToken struct {
}

func (DefToken) String() string {
	return "<Def>"
}

type DefTokenizer struct{}

func (DefTokenizer) Tokenize(word string) (Token, error) {
	if word == "def" {
		return DefToken{}, nil
	} else {
		return nil, ErrTokenNotMatched{Word: word}
	}
}
