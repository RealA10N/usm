package parse_test

import (
	"testing"
	"usm/lex"
	"usm/parse"
	"usm/source"

	"github.com/stretchr/testify/assert"
)

func TestInstructionParserMultipleTargets(t *testing.T) {
	v, ctx := source.NewSourceView("%div %mod = divmod %x %y").Detach()
	t1 := lex.Token{Type: lex.RegToken, View: v.Subview(0, 4)}
	t2 := lex.Token{Type: lex.RegToken, View: v.Subview(5, 9)}
	eq := lex.Token{Type: lex.EqlToken, View: v.Subview(10, 11)}
	op := lex.Token{Type: lex.OprToken, View: v.Subview(12, 18)}
	a1 := lex.Token{Type: lex.RegToken, View: v.Subview(19, 21)}
	a2 := lex.Token{Type: lex.RegToken, View: v.Subview(22, 24)}
	tknView := parse.NewTokenView([]lex.Token{
		t1, t2, eq, op, a1, a2,
	})

	expected := parse.InstructionNode{
		Operator: v.Subview(12, 18),
		Arguments: []parse.ArgumentNode{
			parse.ArgumentNode{v.Subview(19, 21)},
			parse.ArgumentNode{v.Subview(22, 24)},
		},
		Targets: []parse.RegisterNode{
			parse.RegisterNode{v.Subview(0, 4)},
			parse.RegisterNode{v.Subview(5, 9)},
		},
	}

	inst, err := parse.InstructionParser{}.Parse(&tknView)
	assert.Nil(t, err)
	assert.Equal(t, expected, inst)
	assert.Equal(t, v, inst.View())
	assert.Equal(t, "%div %mod = divmod %x %y", inst.String(ctx))
}
