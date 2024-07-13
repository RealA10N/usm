package tokens

type LcurlToken struct {
}

func (LcurlToken) String() string {
	return "<Lcurl>"
}

type LcurlTokenizer struct {
}

func (LcurlTokenizer) Tokenize(word string) (Token, error) {
	if word == "{" {
		return LcurlToken{}, nil
	} else {
		return nil, ErrTokenNotMatched{word}
	}
}
