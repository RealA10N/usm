package lex

import (
	"errors"
	"unicode"
)

type OprTokenizer struct{}

func (OprTokenizer) Tokenize(txt *SourceView) (tkn Token, err error) {
	idx := txt.Index(unicode.IsSpace)
	if idx == 0 {
		err = errors.New("token not matched")
		return
	}

	tknView, restView := txt.Partition(idx)
	detachedTkn, _ := tknView.Detach()
	tkn = Token{Type: OprToken, View: detachedTkn}
	*txt = restView
	return
}
