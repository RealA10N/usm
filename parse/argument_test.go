package parse_test

import (
	"testing"
	"usm/lex"
	"usm/parse"
	"usm/source"

	"github.com/RealA10N/view"
	"github.com/stretchr/testify/assert"
)

func TestArgumentNodeParserSimpleCase(t *testing.T) {
	p := parse.ArgumentNodeParser{}
	v, ctx := source.NewSourceView("$i32 %0 .entry").Detach()

	typView := v.Subview(0, 4)
	regView := v.Subview(5, 7)
	lblView := v.Subview(8, 14)
	typTkn := lex.Token{Type: lex.TypToken, View: typView}
	regTkn := lex.Token{Type: lex.RegToken, View: regView}
	lblTkn := lex.Token{Type: lex.LblToken, View: lblView}
	tkns := view.NewView[lex.Token, uint32]([]lex.Token{typTkn, regTkn, lblTkn})
	expectedSubview := tkns.Subview(2, 3)

	node, err := p.Parse(&tkns)
	assert.Nil(t, err)
	assert.Equal(t, expectedSubview, tkns)

	assert.EqualValues(t, 0, node.View.Start)
	assert.EqualValues(t, 7, node.View.End)
	assert.Equal(t, "$i32 %0", string(node.View.Raw(ctx)))

	assert.Equal(t, typTkn, node.Type)
	assert.Equal(t, regTkn, node.Register)
}