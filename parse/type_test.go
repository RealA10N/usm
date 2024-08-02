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
	p := parse.TypeParser{}
	typView, ctx := source.NewSourceView("$i32").Detach()
	typTkn := lex.Token{Type: lex.TypToken, View: typView}
	tkns := view.NewView[lex.Token, uint32]([]lex.Token{typTkn})
	expectedSubview := tkns.Subview(1, 1)

	node, err := p.Parse(&tkns)
	assert.Nil(t, err)
	assert.Equal(t, expectedSubview, tkns)

	assert.Equal(t, "$i32", string(node.View.Raw(ctx)))
	assert.EqualValues(t, 0, node.View.Start)
	assert.EqualValues(t, 4, node.View.End)
}
