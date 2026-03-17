package parse_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

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

// TestFunctionWithVariableRoundTrip verifies that a function using variables
// in its instructions round-trips through String() correctly.
func TestFunctionWithVariableRoundTrip(t *testing.T) {
	src := "func @foo {\n\tstore &x %n\n\tret\n}"

	srcView := core.NewSourceView(src)
	result, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	v := parse.NewTokenView(result)
	function, perr := parse.NewFunctionParser().Parse(&v)
	assert.Nil(t, perr)

	assert.NotNil(t, function.Instructions)
	assert.Len(t, function.Instructions.Nodes, 2)

	strCtx := parse.StringContext{SourceContext: srcView.Ctx()}
	assert.Equal(t, src, function.String(&strCtx))
}
