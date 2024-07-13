package tokens

import (
	"fmt"
	"strings"
)

type TypeToken struct {
	name string
}

func (tkn TypeToken) String() string {
	return fmt.Sprintf("<Type %v>", tkn.name)
}

type TypeTokenizer struct {
}

func (TypeTokenizer) Tokenize(word string) (Token, error) {
	name, ok := strings.CutPrefix(word, "$")
	if ok {
		return TypeToken{name}, nil
	} else {
		return nil, ErrTokenNotMatched{word}
	}
}
