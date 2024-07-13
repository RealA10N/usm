package tokens

import (
	"fmt"
	"strings"
)

type RegisterToken struct {
	name string
}

func (token RegisterToken) String() string {
	return fmt.Sprintf("<Register %v>", token.name)
}

type RegisterTokenizer struct {
}

func (RegisterTokenizer) Tokenize(word string) (Token, error) {
	name, ok := strings.CutPrefix(word, "%")
	if ok {
		return RegisterToken{name}, nil
	} else {
		return nil, ErrTokenNotMatched{word}
	}
}
