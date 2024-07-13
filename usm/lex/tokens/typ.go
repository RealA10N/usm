package tokens

import (
	"fmt"
	"strings"
	"usm/lex/base"
)

type TypToken struct {
	Name string
}

func (token TypToken) String() string {
	return fmt.Sprintf("<Typ %v>", token.Name)
}

type TypTokenizer struct {
}

func (TypTokenizer) Tokenize(word string) (base.Token, error) {
	name, ok := strings.CutPrefix(word, "$")
	if ok {
		return TypToken{name}, nil
	} else {
		return nil, base.ErrTokenNotMatched{Word: word}
	}
}
