package tokens

import (
	"fmt"
	"strings"
	"usm/lex/base"
)

type OprToken struct {
	Name string
}

func (token OprToken) String() string {
	return fmt.Sprintf("<Opr %v>", token.Name)
}

type OprTokenizer struct {
}

func (OprTokenizer) Tokenize(word string) (base.Token, error) {
	return OprToken{Name: strings.ToLower(word)}, nil
}
