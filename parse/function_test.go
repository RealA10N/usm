package parse_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func TestFunctionNoBody(t *testing.T) {
	src := "func @foo"

	expected := parse.FunctionNode{
		UnmanagedSourceView: core.UnmanagedSourceView{Start: 0, End: 9},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          core.UnmanagedSourceView{Start: 5, End: 9},
		},
	}

	testExpectedFunctionParsing(t, src, expected, src)
}

func TestFunctionOneLineZeroInstructions(t *testing.T) {
	src := "func @foo { }"

	expected := parse.FunctionNode{
		UnmanagedSourceView: core.UnmanagedSourceView{Start: 0, End: 13},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          core.UnmanagedSourceView{Start: 5, End: 9},
		},
		Instructions: &parse.BlockNode[parse.InstructionNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 10, End: 13},
		},
	}

	testExpectedFunctionParsing(t, src, expected, src)
}

func TestFunctionOneLine(t *testing.T) {
	src := "func @foo { $32 %0 = bar\n}"
	expected := parse.FunctionNode{
		UnmanagedSourceView: core.UnmanagedSourceView{Start: 0, End: 26},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          core.UnmanagedSourceView{Start: 5, End: 9},
		},
		Instructions: &parse.BlockNode[parse.InstructionNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 10, End: 26},
			Nodes: []parse.InstructionNode{
				{
					Operator: core.UnmanagedSourceView{Start: 21, End: 24},
					Targets: []parse.TargetNode{
						{
							Type:     &parse.TypeNode{Identifier: core.UnmanagedSourceView{Start: 12, End: 15}},
							Register: parse.RegisterNode{TokenNode: parse.TokenNode{core.UnmanagedSourceView{Start: 16, End: 18}}},
						},
					},
				},
			},
		},
	}

	expectedString := "func @foo {\n\t$32 %0 = bar\n}"

	testExpectedFunctionParsing(t, src, expected, expectedString)
}

// TestFunctionEmptyBodyWithTrailingComment verifies that a block with no
// instructions but a whole-line trailing comment is NOT collapsed to "{ }".
// The comment must be rendered indented inside the braces.
func TestFunctionEmptyBodyWithTrailingComment(t *testing.T) {
	src := "func @foo {\n\t; trailing comment\n}"

	expected := parse.FunctionNode{
		UnmanagedSourceView: core.UnmanagedSourceView{Start: 0, End: 33},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          core.UnmanagedSourceView{Start: 5, End: 9},
		},
		Instructions: &parse.BlockNode[parse.InstructionNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 10, End: 33},
			TrailingComments: []lex.Comment{
				{View: core.UnmanagedSourceView{Start: 13, End: 31}},
			},
		},
	}

	expectedString := "func @foo {\n\t; trailing comment\n}"

	testExpectedFunctionParsing(t, src, expected, expectedString)
}

// TestFunctionWithLeadingCommentOnInstruction verifies that a whole-line comment
// before an instruction is attached to that instruction's LeadingComments field.
func TestFunctionWithLeadingCommentOnInstruction(t *testing.T) {
	src := "func @foo {\n\t; step one\n\tret\n}"
	// byte offsets: { at 10, comment [13,23), ret [25,28), } at 29
	expected := parse.FunctionNode{
		UnmanagedSourceView: core.UnmanagedSourceView{Start: 0, End: 30},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          core.UnmanagedSourceView{Start: 5, End: 9},
		},
		Instructions: &parse.BlockNode[parse.InstructionNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 10, End: 30},
			Nodes: []parse.InstructionNode{
				{
					Operator:        core.UnmanagedSourceView{Start: 25, End: 28},
					LeadingComments: []lex.Comment{{View: core.UnmanagedSourceView{Start: 13, End: 23}}},
				},
			},
		},
	}

	testExpectedFunctionParsing(t, src, expected, src)
}

// TestFunctionWithTrailingCommentOnInstruction verifies that an inline comment
// after an instruction's last token is attached to TrailingComment.
func TestFunctionWithTrailingCommentOnInstruction(t *testing.T) {
	src := "func @foo {\n\tret ; done\n}"
	// byte offsets: { at 10, ret [13,16), comment [17,23), } at 24
	comment := lex.Comment{View: core.UnmanagedSourceView{Start: 17, End: 23}}

	expected := parse.FunctionNode{
		UnmanagedSourceView: core.UnmanagedSourceView{Start: 0, End: 25},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          core.UnmanagedSourceView{Start: 5, End: 9},
		},
		Instructions: &parse.BlockNode[parse.InstructionNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 10, End: 25},
			Nodes: []parse.InstructionNode{
				{
					Operator:        core.UnmanagedSourceView{Start: 13, End: 16},
					TrailingComment: &comment,
				},
			},
		},
	}

	testExpectedFunctionParsing(t, src, expected, src)
}

// TestFunctionBodyWithTrailingComment verifies that a whole-line comment after
// the last instruction (before '}') ends up in BlockNode.TrailingComments.
func TestFunctionBodyWithTrailingComment(t *testing.T) {
	src := "func @foo {\n\tret\n\t; the end\n}"
	// byte offsets: { at 10, ret [13,16), comment [18,27), } at 28
	expected := parse.FunctionNode{
		UnmanagedSourceView: core.UnmanagedSourceView{Start: 0, End: 29},
		Signature: parse.FunctionSignatureNode{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 5, End: 9},
			Identifier:          core.UnmanagedSourceView{Start: 5, End: 9},
		},
		Instructions: &parse.BlockNode[parse.InstructionNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{Start: 10, End: 29},
			Nodes: []parse.InstructionNode{
				{Operator: core.UnmanagedSourceView{Start: 13, End: 16}},
			},
			TrailingComments: []lex.Comment{{View: core.UnmanagedSourceView{Start: 18, End: 27}}},
		},
	}

	testExpectedFunctionParsing(t, src, expected, src)
}

// MARK: Helpers

func testExpectedFunctionParsing(
	t *testing.T,
	src string,
	expectedStructure parse.FunctionNode,
	expectedString string,
) {
	t.Helper()

	srcView := core.NewSourceView(src)
	result, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	v := parse.NewTokenView(result)
	function, perr := parse.NewFunctionParser().Parse(&v)
	assert.Nil(t, perr)

	assert.Equal(t, expectedStructure, function)
	strCtx := parse.StringContext{SourceContext: srcView.Ctx()}
	assert.Equal(t, expectedString, function.String(&strCtx))
}
