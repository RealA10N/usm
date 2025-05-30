package lex

import (
	"errors"
	"unicode"

	"alon.kr/x/usm/core"
)

// Tokenizer that scans words prefixed with the provided prefix string,
// Scans the word until a whitespace is encountered.
type PrefixedTokenizer struct {
	Prefix string
	Token  TokenType
}

func (t PrefixedTokenizer) Tokenize(txt *core.SourceView) (tkn Token, err error) {
	ok := txt.HasPrefix(core.NewSourceView(t.Prefix))
	if !ok {
		err = errors.New("token not matched")
		return
	}

	idx := txt.IndexFunc(unicode.IsSpace)
	tknView, restView := txt.Partition(idx)
	detachedTkn, _ := tknView.Detach()
	tkn = Token{Type: t.Token, View: detachedTkn}
	*txt = restView
	return
}
