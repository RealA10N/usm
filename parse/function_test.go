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
	src := "func @foo { $32 %0 = bar }"
	expected := parse.FunctionNode{
		UnmanagedSourceView: source.UnmanagedSourceView{Start: 0, End: 26},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: source.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          source.UnmanagedSourceView{Start: 5, End: 9},
		},
		Instructions: &parse.BlockNode[parse.InstructionNode]{
			UnmanagedSourceView: source.UnmanagedSourceView{Start: 10, End: 26},
			Nodes: []parse.InstructionNode{
				{
					Operator: source.UnmanagedSourceView{Start: 21, End: 24},
					Targets: []parse.ParameterNode{
						{
							Type:     parse.TypeNode{Identifier: source.UnmanagedSourceView{Start: 12, End: 15}},
							Register: parse.RegisterNode{source.UnmanagedSourceView{Start: 16, End: 18}},
						},
					},
				},
			},
		},
	}

	expectedString := "func @foo {\n\t$32 %0 = bar\n}"

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
	tkns, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	v := parse.NewTokenView(tkns)
	function, perr := parse.NewFunctionParser().Parse(&v)
	assert.Nil(t, perr)

	assert.Equal(t, expectedStructure, function)
	strCtx := parse.StringContext{SourceContext: srcView.Ctx()}
	assert.Equal(t, expectedString, function.String(&strCtx))
}
