package tokens

import (
	"fmt"
	"strings"
)

type GlobalToken struct {
	name string
}

func (tkn GlobalToken) String() string {
	return fmt.Sprintf("<Global %v>", tkn.name)
}

type GlobalTokenizer struct {
}

func (GlobalTokenizer) Tokenize(word string) (Token, error) {
	name, ok := strings.CutPrefix(word, "@")
	if ok {
		return RegisterToken{name}, nil
	} else {
		return nil, ErrTokenNotMatched{word}
	}
}
