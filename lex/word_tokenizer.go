package lex

import (
	"errors"
	"unicode"

	"alon.kr/x/usm/source"
)

type WordTokenizer struct {
	Token TokenType
}

func (t WordTokenizer) Tokenize(txt *source.SourceView) (tkn Token, err error) {
	idx := txt.Index(unicode.IsSpace) // TODO: currently assuming prefix has no whitespaces.
	if idx < 1 {
		err = errors.New("token not matched")
		return
	}

	tknView, restView := txt.Partition(idx)
	detachedTkn, _ := tknView.Detach()
	tkn = Token{Type: t.Token, View: detachedTkn}
	*txt = restView
	return
}
