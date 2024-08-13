package parse_test

import (
	"testing"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"

	"github.com/stretchr/testify/assert"
)

func TestInstructionParserMultipleTargets(t *testing.T) {
	srcView := source.NewSourceView("%div %mod = divmod %x %y")
	unmanaged := srcView.Unmanaged()

	expected := parse.InstructionNode{
		Operator: unmanaged.Subview(12, 18),
		Arguments: []parse.ArgumentNode{
			parse.RegisterNode{unmanaged.Subview(19, 21)},
			parse.RegisterNode{unmanaged.Subview(22, 24)},
		},
		Targets: []parse.RegisterNode{
			{unmanaged.Subview(0, 4)},
			{unmanaged.Subview(5, 9)},
		},
	}

	expectedString := "\t%div %mod = divmod %x %y\n"

	testExpectedInstruction(t, srcView, expected, expectedString)
}

func TestInstructionWithImmediateValuesAndLabel(t *testing.T) {
	srcView := source.NewSourceView(".entry %res = add %x $32 #1 .arg")
	unmanaged := srcView.Unmanaged()

	expected := parse.InstructionNode{
		Operator: unmanaged.Subview(14, 17),
		Arguments: []parse.ArgumentNode{
			parse.RegisterNode{unmanaged.Subview(18, 20)},
			parse.ImmediateNode{
				Type: parse.TypeNode{Identifier: unmanaged.Subview(21, 24)},
				Value: parse.ImmediateFinalValueNode{
					unmanaged.Subview(25, 27),
				},
			},
			parse.LabelNode{unmanaged.Subview(28, 32)},
		},
		Targets: []parse.RegisterNode{
			{unmanaged.Subview(7, 11)},
		},
		Labels: []parse.LabelNode{
			{unmanaged.Subview(0, 6)},
		},
	}

	expectedString := ".entry\n\t%res = add %x $32 #1 .arg\n"

	testExpectedInstruction(t, srcView, expected, expectedString)
}

// MARK: Helpers

func testExpectedInstruction(
	t *testing.T,
	srcView source.SourceView,
	expected parse.InstructionNode,
	expectedString string,
) {
	t.Helper()

	tkns, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	tknView := parse.NewTokenView(tkns)
	inst, perr := parse.InstructionParser{}.Parse(&tknView)
	assert.Nil(t, perr)

	assert.Equal(t, expected, inst)
	strCtx := parse.StringContext{SourceContext: srcView.Ctx()}
	assert.Equal(t, expectedString, inst.String(&strCtx))
}
