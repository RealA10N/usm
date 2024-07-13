package tokens

import (
	"fmt"
	"strings"
	"usm/lex/base"
)

type OprToken struct {
	name string
}

func (token OprToken) String() string {
	return fmt.Sprintf("<Opr %v>", token.name)
}

type OprTokenizer struct {
}

func (OprTokenizer) Tokenize(word string) (base.Token, error) {
	return OprToken{name: strings.ToLower(word)}, nil
}
