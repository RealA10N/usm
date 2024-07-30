package tokens

import (
	"fmt"
	"strings"
	"usm/lex/base"
)

type ImmToken struct {
	Value string
}

func (token ImmToken) String() string {
	return fmt.Sprintf("<Imm #%v>", token.Value)
}

type ImmTokenizer struct{}

func (ImmTokenizer) Tokenize(word string) (base.Token, error) {
	value, ok := strings.CutPrefix(word, "#")
	if ok {
		return ImmToken{Value: value}, nil
	} else {
		return nil, base.ErrTokenNotMatched{Word: word}
	}
}
