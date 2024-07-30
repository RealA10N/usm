package tokens

import (
	"fmt"
	"strings"
	"usm/lex/base"
)

type LblToken struct {
	Name string
}

func (token LblToken) String() string {
	return fmt.Sprintf("<Lbl .%v>", token.Name)
}

type LblTokenizer struct {
}

func (LblTokenizer) Tokenize(word string) (base.Token, error) {
	name, ok := strings.CutPrefix(word, ".")
	if ok {
		return LblToken{name}, nil
	} else {
		return nil, base.ErrTokenNotMatched{Word: word}
	}
}
