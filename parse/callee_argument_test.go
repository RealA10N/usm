package parse_test

import (
	"testing"
	"usm/lex"
	"usm/parse"
	"usm/source"

	"github.com/RealA10N/view"
	"github.com/stretchr/testify/assert"
)

func TestArgumentParserSimpleCase(t *testing.T) {
	v, ctx := source.NewSourceView("$i32 %0 .entry").Detach()

	typView := v.Subview(0, 4)
	regView := v.Subview(5, 7)
	lblView := v.Subview(8, 14)
	typTkn := lex.Token{Type: lex.TypToken, View: typView}
	regTkn := lex.Token{Type: lex.RegToken, View: regView}
	lblTkn := lex.Token{Type: lex.LblToken, View: lblView}
	tkns := view.NewView[lex.Token, uint32]([]lex.Token{typTkn, regTkn, lblTkn})
	expectedSubview := tkns.Subview(2, 3)

	node, err := parse.CalleeArgumentParser{}.Parse(&tkns)
	assert.Nil(t, err)
	assert.Equal(t, expectedSubview, tkns)

	assert.EqualValues(t, 0, node.View().Start)
	assert.EqualValues(t, 7, node.View().End)
	assert.Equal(t, "$i32 %0", string(node.View().Raw(ctx)))

	assert.Equal(t, "$i32", string(node.Type.View().Raw(ctx)))
	assert.Equal(t, "%0", string(node.Register.View().Raw(ctx)))
}

func TestArgumentTypEofError(t *testing.T) {
	_, ctx := source.NewSourceView("").Detach()
	tkns := []lex.Token{}
	view := view.NewView[lex.Token, uint32](tkns)

	_, err := parse.CalleeArgumentParser{}.Parse(&view)
	assert.NotNil(t, err)
	assert.EqualValues(t, 0, view.Len())
	assert.EqualValues(t, "reached end of file (expected <Type>)", err.Error(ctx))
}

func TestArgumentRegEofError(t *testing.T) {
	v, ctx := source.NewSourceView("$i32").Detach()
	tkn := lex.Token{Type: lex.TypToken, View: v}
	view := view.NewView[lex.Token, uint32]([]lex.Token{tkn})

	_, err := parse.CalleeArgumentParser{}.Parse(&view)
	assert.NotNil(t, err)
	assert.EqualValues(t, 0, view.Len())
	assert.EqualValues(t, "reached end of file (expected <Register>)", err.Error(ctx))
}
