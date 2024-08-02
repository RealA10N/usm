package parse_test

import (
	"testing"
	"usm/lex"
	"usm/parse"
	"usm/source"

	"github.com/RealA10N/view"
	"github.com/stretchr/testify/assert"
)

func TestTypeParserSimpleCase(t *testing.T) {
	typView, ctx := source.NewSourceView("$i32").Detach()
	typTkn := lex.Token{Type: lex.TypToken, View: typView}
	tkns := view.NewView[lex.Token, uint32]([]lex.Token{typTkn})
	expectedSubview := tkns.Subview(1, 1)

	node, err := parse.TypeParser{}.Parse(&tkns)
	assert.Nil(t, err)
	assert.Equal(t, expectedSubview, tkns)

	assert.Equal(t, "$i32", string(node.View.Raw(ctx)))
	assert.EqualValues(t, 0, node.View.Start)
	assert.EqualValues(t, 4, node.View.End)
}

func TestTypeParserEofError(t *testing.T) {
	_, ctx := source.NewSourceView("").Detach()
	tkns := []lex.Token{}
	view := view.NewView[lex.Token, uint32](tkns)

	_, err := parse.TypeParser{}.Parse(&view)
	assert.NotNil(t, err)
	assert.EqualValues(t, 0, view.Len())
	assert.EqualValues(t, "expected <Type> token, but file ended", err.Error(ctx))
}

func TestTypeParserUnexpectedTokenError(t *testing.T) {
	regView, ctx := source.NewSourceView("%0").Detach()
	regTkn := lex.Token{Type: lex.RegToken, View: regView}
	tkns := view.NewView[lex.Token, uint32]([]lex.Token{regTkn})

	_, err := parse.TypeParser{}.Parse(&tkns)
	assert.NotNil(t, err)
	assert.EqualValues(t, "expected <Type> token, but got <Register \"%0\">", err.Error(ctx))
}
