package lex

import "errors"

type RcrTokenizer struct{}

func (RcrTokenizer) Tokenize(txt *SourceView) (tkn Token, err error) {
	chr, err := txt.At(0)
	if err != nil || chr != '}' {
		err = errors.New("token not matched")
		return
	}

	detachedTkn, _ := txt.Subview(0, 1).Detach()
	tkn = Token{Type: RcrToken, View: detachedTkn}
	*txt = txt.Subview(1, txt.Len())
	return
}
