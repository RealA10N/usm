package tokens

type RcurlToken struct {
}

func (RcurlToken) String() string {
	return "<Rcurl>"
}

type RcurlTokenizer struct {
}

func (RcurlTokenizer) Tokenize(word string) (Token, error) {
	if word == "}" {
		return RcurlToken{}, nil
	} else {
		return nil, ErrTokenNotMatched{word}
	}
}
