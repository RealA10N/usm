package parse_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

// TestFunctionWithVariableDeclaration checks that a variable declaration
// preamble is parsed into FunctionNode.Variables, not Instructions.
func TestFunctionWithVariableDeclaration(t *testing.T) {
	src := "func @foo {\n\t&local $32\n\tret\n}"

	srcView := core.NewSourceView(src)
	result, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	v := parse.NewTokenView(result)
	function, perr := parse.NewFunctionParser().Parse(&v)
	assert.Nil(t, perr)

	// One variable declaration, one instruction.
	assert.Len(t, function.Variables, 1)
	assert.NotNil(t, function.Instructions)
	assert.Len(t, function.Instructions.Nodes, 1)

	// Variable name is "&local", type is "$32".
	varDecl := function.Variables[0]
	strCtx := parse.StringContext{SourceContext: srcView.Ctx()}
	assert.Equal(t, "&local", varDecl.Variable.String(&strCtx))
	assert.Equal(t, "$32", varDecl.Type.String(&strCtx))

	// Round-trip: String() must reproduce the original source.
	assert.Equal(t, src, function.String(&strCtx))
}

// TestFunctionWithMultipleVariableDeclarations checks that multiple variable
// declarations are all captured and that instructions follow normally.
func TestFunctionWithMultipleVariableDeclarations(t *testing.T) {
	src := "func @foo {\n\t&a $32\n\t&b $32\n\tret\n}"

	srcView := core.NewSourceView(src)
	result, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	v := parse.NewTokenView(result)
	function, perr := parse.NewFunctionParser().Parse(&v)
	assert.Nil(t, perr)

	assert.Len(t, function.Variables, 2)
	assert.Len(t, function.Instructions.Nodes, 1)
}

// TestVariableDeclarationWithTrailingComment checks that inline comments on
// variable declaration lines are captured and preserved.
func TestVariableDeclarationWithTrailingComment(t *testing.T) {
	src := "func @foo {\n\t&x $32 ; a counter\n\tret\n}"

	srcView := core.NewSourceView(src)
	result, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	v := parse.NewTokenView(result)
	function, perr := parse.NewFunctionParser().Parse(&v)
	assert.Nil(t, perr)

	assert.Len(t, function.Variables, 1)
	assert.NotNil(t, function.Variables[0].TrailingComment)

	strCtx := parse.StringContext{SourceContext: srcView.Ctx()}
	assert.Equal(t, src, function.String(&strCtx))
}

// TestFunctionWithNoVariables verifies that a function without variable
// declarations still works as before.
func TestFunctionWithNoVariables(t *testing.T) {
	src := "func @foo {\n\tret\n}"

	srcView := core.NewSourceView(src)
	result, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	v := parse.NewTokenView(result)
	function, perr := parse.NewFunctionParser().Parse(&v)
	assert.Nil(t, perr)

	assert.Empty(t, function.Variables)
	assert.Len(t, function.Instructions.Nodes, 1)
}

// TestVariableAsArgument verifies that &var can appear as an instruction
// argument and is tokenized as a VariableToken.
func TestVariableAsArgument(t *testing.T) {
	src := "store &local %x\n"

	srcView := core.NewSourceView(src)
	result, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	// Check the token stream contains a VariableToken.
	found := false
	for _, tkn := range result {
		if tkn.Type == lex.VariableToken {
			found = true
			assert.Equal(t, "&local", string(tkn.View.Raw(srcView.Ctx())))
		}
	}
	assert.True(t, found, "expected VariableToken in token stream")

	v := parse.NewTokenView(result)
	node, perr := parse.NewInstructionParser().Parse(&v)
	assert.Nil(t, perr)

	// The instruction has no targets, one variable arg and one register arg.
	assert.Len(t, node.Targets, 0)
	assert.Len(t, node.Arguments, 2)
	_, isVar := node.Arguments[0].(parse.VariableNode)
	assert.True(t, isVar)
}
