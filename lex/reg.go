package lex

import (
	"errors"
	"unicode"
)

type RegTokenizer struct{}

func (RegTokenizer) Tokenize(txt *SourceView) (tkn Token, err error) {
	chr, err := txt.At(0)
	if err != nil || chr != '%' {
		err = errors.New("token not matched")
		return
	}

	idx := txt.Index(unicode.IsSpace)
	tknView, restView := txt.Partition(idx)
	detachedTkn, _ := tknView.Detach()
	tkn = Token{Type: RegToken, View: detachedTkn}
	*txt = restView
	return
}
