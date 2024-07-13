package tokens

import (
	"fmt"
	"strings"
	"usm/lex/base"
)

type RegToken struct {
	name string
}

func (token RegToken) String() string {
	return fmt.Sprintf("<Reg %v>", token.name)
}

type RegTokenizer struct {
}

func (RegTokenizer) Tokenize(word string) (base.Token, error) {
	name, ok := strings.CutPrefix(word, "%")
	if ok {
		return RegToken{name}, nil
	} else {
		return nil, base.ErrTokenNotMatched{Word: word}
	}
}
