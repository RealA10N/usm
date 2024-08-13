package parse_test

import (
	"testing"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"

	"github.com/stretchr/testify/assert"
)

func TestParameterParserSimpleCase(t *testing.T) {
	v, ctx := source.NewSourceView("$i32 %0 .entry").Detach()

	typView := v.Subview(0, 4)
	regView := v.Subview(5, 7)
	lblView := v.Subview(8, 14)
	typTkn := lex.Token{Type: lex.TypeToken, View: typView}
	regTkn := lex.Token{Type: lex.RegisterToken, View: regView}
	lblTkn := lex.Token{Type: lex.LabelToken, View: lblView}
	tkns := parse.NewTokenView([]lex.Token{typTkn, regTkn, lblTkn})
	expectedSubview := parse.TokenView{tkns.Subview(2, 3)}

	node, err := parse.ParameterParser{}.Parse(&tkns)
	assert.Nil(t, err)
	assert.Equal(t, expectedSubview, tkns)

	assert.EqualValues(t, 0, node.View().Start)
	assert.EqualValues(t, 7, node.View().End)
	assert.Equal(t, "$i32 %0", string(node.View().Raw(ctx)))

	assert.Equal(t, "$i32", string(node.Type.View().Raw(ctx)))
	assert.Equal(t, "%0", string(node.Register.View().Raw(ctx)))
}

func TestParameterTypEofError(t *testing.T) {
	view := parse.NewTokenView([]lex.Token{})
	_, err := parse.ParameterParser{}.Parse(&view)
	expected := parse.EofError{Expected: []lex.TokenType{lex.TypeToken}}
	assert.Equal(t, expected, err)
}

func TestParameterRegEofError(t *testing.T) {
	v := source.NewSourceView("$i32").Unmanaged()
	tkn := lex.Token{Type: lex.TypeToken, View: v}
	view := parse.NewTokenView([]lex.Token{tkn})

	expected := parse.UnexpectedTokenError{
		Expected: []lex.TokenType{lex.RegisterToken},
		Actual:   tkn,
	}

	_, err := parse.ParameterParser{}.Parse(&view)
	assert.Equal(t, expected, err)
}
