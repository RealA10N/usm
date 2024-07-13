package tokens

import (
	"fmt"
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
	return OprToken{name: word}, nil
}
