package lex

import (
	"errors"
	"usm/source"
)

type KeywordTokenizer struct {
	Keyword string
	Token   TokenType
}

func (t KeywordTokenizer) Tokenize(txt *source.SourceView) (tkn Token, err error) {
	keywordView := source.NewSourceView(t.Keyword)
	ok := txt.HasPrefix(keywordView)
	if !ok {
		err = errors.New("token not matched")
		return
	}

	tknView, restView := txt.Partition(keywordView.Len())
	unmanagedTknView, _ := tknView.Detach()
	tkn = Token{Type: t.Token, View: unmanagedTknView}
	*txt = restView
	return
}
