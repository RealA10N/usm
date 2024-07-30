package tokens

import (
	"fmt"
	"strings"
	"usm/lex/base"
)

type GlbToken struct {
	Name string
}

func (token GlbToken) String() string {
	return fmt.Sprintf("<Glb @%v>", token.Name)
}

type GlbTokenizer struct{}

func (GlbTokenizer) Tokenize(word string) (base.Token, error) {
	name, ok := strings.CutPrefix(word, "@")
	if ok {
		return GlbToken{name}, nil
	} else {
		return nil, base.ErrTokenNotMatched{Word: word}
	}
}
