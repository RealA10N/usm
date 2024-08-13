package parse_test

import (
	"testing"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"
	"github.com/stretchr/testify/assert"
)

func TestFunctionNoBody(t *testing.T) {
	src := "func @foo"

	expected := parse.FunctionNode{
		UnmanagedSourceView: source.UnmanagedSourceView{Start: 0, End: 9},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: source.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          source.UnmanagedSourceView{Start: 5, End: 9},
		},
	}

	testExpectedFunctionParsing(t, src, expected, src)
}

func TestFunctionOneLineZeroInstructions(t *testing.T) {
	src := "func @foo { }"

	expected := parse.FunctionNode{
		UnmanagedSourceView: source.UnmanagedSourceView{Start: 0, End: 13},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: source.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          source.UnmanagedSourceView{Start: 5, End: 9},
		},
		Instructions: &parse.BlockNode[parse.InstructionNode]{
			UnmanagedSourceView: source.UnmanagedSourceView{Start: 10, End: 13},
		},
	}

	testExpectedFunctionParsing(t, src, expected, src)
}

func TestFunctionOneLine(t *testing.T) {
	src := "func @foo { %0 = bar }"
	expected := parse.FunctionNode{
		UnmanagedSourceView: source.UnmanagedSourceView{Start: 0, End: 22},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: source.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          source.UnmanagedSourceView{Start: 5, End: 9},
		},
		Instructions: &parse.BlockNode[parse.InstructionNode]{
			UnmanagedSourceView: source.UnmanagedSourceView{Start: 10, End: 22},
			Nodes: []parse.InstructionNode{
				{
					Operator: source.UnmanagedSourceView{Start: 17, End: 20},
					Targets: []parse.RegisterNode{
						{source.UnmanagedSourceView{Start: 12, End: 14}},
					},
				},
			},
		},
	}

	expectedString := "func @foo {\n\t%0 = bar\n}"

	testExpectedFunctionParsing(t, src, expected, expectedString)
}

// MARK: Helpers

func testExpectedFunctionParsing(
	t *testing.T,
	src string,
	expectedStructure parse.FunctionNode,
	expectedString string,
) {
	t.Helper()

	srcView := source.NewSourceView(src)
	ctx := source.SourceContext{ViewContext: srcView.Ctx()}
	tkns, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	v := parse.NewTokenView(tkns)
	function, perr := parse.NewFunctionParser().Parse(&v)
	assert.Nil(t, perr)

	assert.Equal(t, expectedStructure, function)
	assert.Equal(t, expectedString, function.String(ctx))
}
