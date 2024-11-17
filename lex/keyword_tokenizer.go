package lex

import (
	"errors"

	"alon.kr/x/usm/core"
)

type KeywordTokenizer struct {
	Keyword string
	Token   TokenType
}

func (t KeywordTokenizer) Tokenize(txt *core.SourceView) (tkn Token, err error) {
	keywordView := core.NewSourceView(t.Keyword)
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
