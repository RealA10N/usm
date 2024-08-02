package lex

import (
	"errors"
	"unicode"
	"usm/source"
)

type ImmTokenizer struct{}

func (ImmTokenizer) Tokenize(txt *source.SourceView) (tkn Token, err error) {
	chr, err := txt.At(0)
	if err != nil || chr != '#' {
		err = errors.New("token not matched")
		return
	}

	idx := txt.Index(unicode.IsSpace)
	tknView, restView := txt.Partition(idx)
	detachedTkn, _ := tknView.Detach()
	tkn = Token{Type: ImmToken, View: detachedTkn}
	*txt = restView
	return
}
