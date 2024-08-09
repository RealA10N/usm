package lex

import (
	"errors"
	"unicode"

	"alon.kr/x/usm/source"
)

// Tokenizer that scans words prefixed with the provided prefix string,
// Scans the word until a whitespace is encountered.
type PrefixedTokenizer struct {
	Prefix string
	Token  TokenType
}

func (t PrefixedTokenizer) Tokenize(txt *source.SourceView) (tkn Token, err error) {
	ok := txt.HasPrefix(source.NewSourceView(t.Prefix))
	if !ok {
		err = errors.New("token not matched")
		return
	}

	idx := txt.Index(unicode.IsSpace)
	tknView, restView := txt.Partition(idx)
	detachedTkn, _ := tknView.Detach()
	tkn = Token{Type: t.Token, View: detachedTkn}
	*txt = restView
	return
}
