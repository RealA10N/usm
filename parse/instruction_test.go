package parse_test

import (
	"testing"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"

	"github.com/stretchr/testify/assert"
)

func TestInstructionParserMultipleTargets(t *testing.T) {
	v, ctx := source.NewSourceView("%div %mod = divmod %x %y").Detach()
	t1 := lex.Token{Type: lex.RegisterToken, View: v.Subview(0, 4)}
	t2 := lex.Token{Type: lex.RegisterToken, View: v.Subview(5, 9)}
	eq := lex.Token{Type: lex.EqualToken, View: v.Subview(10, 11)}
	op := lex.Token{Type: lex.OperatorToken, View: v.Subview(12, 18)}
	a1 := lex.Token{Type: lex.RegisterToken, View: v.Subview(19, 21)}
	a2 := lex.Token{Type: lex.RegisterToken, View: v.Subview(22, 24)}
	tknView := parse.NewTokenView([]lex.Token{
		t1, t2, eq, op, a1, a2,
	})

	expected := parse.InstructionNode{
		Operator: v.Subview(12, 18),
		Arguments: []parse.ArgumentNode{
			parse.RegisterNode{v.Subview(19, 21)},
			parse.RegisterNode{v.Subview(22, 24)},
		},
		Targets: []parse.RegisterNode{
			{v.Subview(0, 4)},
			{v.Subview(5, 9)},
		},
	}

	inst, err := parse.InstructionParser{}.Parse(&tknView)
	assert.Nil(t, err)
	assert.Equal(t, expected, inst)
	assert.Equal(t, v, inst.View())
	assert.Equal(t, "\t%div %mod = divmod %x %y\n", inst.String(ctx))
}

func TestInstructionWithImmediateValuesAndLabel(t *testing.T) {
	srcView := source.NewSourceView(".entry %res = add %x $32 #1 .arg")

	expected := parse.InstructionNode{
		Operator: srcView.Unmanaged().Subview(14, 17),
		Arguments: []parse.ArgumentNode{
			parse.RegisterNode{srcView.Unmanaged().Subview(18, 20)},
			parse.ImmediateNode{
				Type:  parse.TypeNode{Identifier: srcView.Unmanaged().Subview(21, 24)},
				Value: srcView.Unmanaged().Subview(25, 27),
			},
			parse.LabelNode{srcView.Unmanaged().Subview(28, 32)},
		},
		Targets: []parse.RegisterNode{
			{srcView.Unmanaged().Subview(7, 11)},
		},
		Labels: []parse.LabelNode{
			{srcView.Unmanaged().Subview(0, 6)},
		},
	}

	tkns, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	tknView := parse.NewTokenView(tkns)
	inst, perr := parse.InstructionParser{}.Parse(&tknView)
	assert.Nil(t, perr)
	assert.Equal(t, expected, inst)

	assert.Equal(t, ".entry\n\t%res = add %x $32 #1 .arg\n", inst.String(srcView.Ctx()))
}
