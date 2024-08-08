package parse_test

import (
	"testing"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"

	"github.com/stretchr/testify/assert"
)

func TestTypeNodeStringer(t *testing.T) {
	typView, ctx := source.NewSourceView("$i32").Detach()
	typTok := lex.Token{Type: lex.TypeToken, View: typView}
	tkns := parse.NewTokenView([]lex.Token{typTok})
	node, err := parse.TypeParser{}.Parse(&tkns)
	assert.Nil(t, err)
	assert.Equal(t, "$i32", node.String(ctx))
}

func TestTypeParserSimpleCase(t *testing.T) {
	typView, ctx := source.NewSourceView("$i32").Detach()
	typTkn := lex.Token{Type: lex.TypeToken, View: typView}
	tkns := parse.NewTokenView([]lex.Token{typTkn})
	expectedSubview := parse.TokenView{tkns.Subview(1, 1)}

	node, err := parse.TypeParser{}.Parse(&tkns)
	assert.Nil(t, err)
	assert.Equal(t, expectedSubview, tkns)

	assert.Equal(t, "$i32", string(node.View().Raw(ctx)))
	assert.EqualValues(t, 0, node.View().Start)
	assert.EqualValues(t, 4, node.View().End)
}

func TestTypeParserEofError(t *testing.T) {
	_, ctx := source.NewSourceView("").Detach()
	tkns := []lex.Token{}
	view := parse.NewTokenView(tkns)

	_, err := parse.TypeParser{}.Parse(&view)
	assert.NotNil(t, err)
	assert.EqualValues(t, 0, view.Len())
	assert.EqualValues(t, "reached end of file (expected <Type>)", err.Error(ctx))
}

func TestTypeParserUnexpectedTokenError(t *testing.T) {
	regView, ctx := source.NewSourceView("%0").Detach()
	regTkn := lex.Token{Type: lex.RegisterToken, View: regView}
	tkns := parse.NewTokenView([]lex.Token{regTkn})

	_, err := parse.TypeParser{}.Parse(&tkns)
	assert.NotNil(t, err)
	assert.EqualValues(t, "got token <Register \"%0\"> (expected <Type>)", err.Error(ctx))
}
