package lex

import (
	"errors"
	"unicode"
	"usm/source"
)

// Tokenizer that scans words prefixed with the provided prefix string,
// Scans the word until a whitespace is encountered.
type prefixedTokenizer struct {
	prefix source.SourceView // TODO: this is not actually a view of a source file, but it is a quick win.
	token  TokenType
}

func NewPrefixedTokenizer(prefix string, token TokenType) prefixedTokenizer {
	return prefixedTokenizer{prefix: source.NewSourceView(prefix), token: token}
}

func (t prefixedTokenizer) Tokenize(txt *source.SourceView) (tkn Token, err error) {
	ok := txt.HasPrefix(t.prefix)
	if !ok {
		err = errors.New("token not matched")
		return
	}

	idx := txt.Index(unicode.IsSpace) // TODO: currently assuming prefix has no whitespaces.
	tknView, restView := txt.Partition(idx)
	detachedTkn, _ := tknView.Detach()
	tkn = Token{Type: t.token, View: detachedTkn}
	*txt = restView
	return
}
