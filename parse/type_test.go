package parse_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"

	"github.com/stretchr/testify/assert"
)

func TestTypeNodeStringer(t *testing.T) {
	typView, ctx := core.NewSourceView("$i32").Detach()
	typTok := lex.Token{Type: lex.TypeToken, View: typView}
	tkns := parse.NewTokenView([]lex.Token{typTok})
	node, err := parse.TypeParser{}.Parse(&tkns)
	assert.Nil(t, err)
	assert.Equal(t, "$i32", node.String(&parse.StringContext{SourceContext: ctx}))
}

func TestTypeParserSimpleCase(t *testing.T) {
	typView, ctx := core.NewSourceView("$i32").Detach()
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
	tkns := []lex.Token{}
	view := parse.NewTokenView(tkns)
	expected := parse.EofError{Expected: []lex.TokenType{lex.TypeToken}}

	_, err := parse.TypeParser{}.Parse(&view)
	assert.Equal(t, expected, err)
}

func TestTypeParserUnexpectedTokenError(t *testing.T) {
	regView := core.NewSourceView("%0").Unmanaged()
	regTkn := lex.Token{Type: lex.RegisterToken, View: regView}
	tkns := parse.NewTokenView([]lex.Token{regTkn})

	expected := parse.UnexpectedTokenError{
		Expected: []lex.TokenType{lex.TypeToken},
		Actual:   regTkn,
	}

	_, err := parse.TypeParser{}.Parse(&tkns)
	assert.Equal(t, expected, err)
}
