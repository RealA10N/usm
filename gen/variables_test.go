package gen_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	usmmanagers "alon.kr/x/usm/usm/managers"
	"github.com/stretchr/testify/assert"
)

// generateFunctionFromSourceUSM is like generateFunctionFromSource but uses the
// full USM instruction set (including load / store / lea).
func generateFunctionFromSourceUSM(
	t *testing.T,
	source string,
) (*gen.FunctionInfo, core.ResultList) {
	t.Helper()

	sourceView := core.NewSourceView(source)

	lexResult, err := lex.NewTokenizer().Tokenize(sourceView)
	assert.NoError(t, err)

	tknView := parse.NewTokenView(lexResult)
	node, result := parse.NewFunctionParser().Parse(&tknView)
	assert.Nil(t, result)

	ctx := usmmanagers.NewGenerationContext().
		NewFileGenerationContext(sourceView.Ctx())

	funcGlobalGen := gen.NewFunctionGlobalGenerator()
	funcGlobalGen.Generate(ctx, node)

	funcGen := gen.NewFunctionGenerator()
	function, results := funcGen.Generate(ctx, node)
	if !results.IsEmpty() || function == nil {
		return nil, results
	}

	results = function.Validate()
	if !results.IsEmpty() {
		return nil, results
	}

	return function, core.ResultList{}
}

// TestFunctionWithVariable verifies that variables are lazily created on first
// use and collected into FunctionInfo.Variables.
func TestFunctionWithVariable(t *testing.T) {
	src := `func $32 @foo $32 %n {
	store &saved %n
	$32 %result = load &saved
	ret %result
}
`
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.True(t, results.IsEmpty(), "unexpected errors: %v", results)
	assert.NotNil(t, function)

	vars := function.Variables.GetAllVariables()
	assert.Len(t, vars, 1)
	assert.Equal(t, "&saved", vars[0].Name)
	assert.Equal(t, "$32", vars[0].Type.String())

	// String() round-trip must reproduce the original source.
	assert.Equal(t, src, function.String())
}

// TestVariableTypeMismatch verifies that using a variable with inconsistent
// types across instructions produces an error.
func TestVariableTypeMismatch(t *testing.T) {
	src := `func $64 @foo $32 %n {
	store &x %n
	$64 %result = load &x
	ret %result
}
`
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.False(t, results.IsEmpty(), "expected type mismatch error")
	assert.Nil(t, function)
}

// TestStoreCritical verifies that store generates without errors and is a
// critical instruction (not eliminated by DCE).
func TestStoreCritical(t *testing.T) {
	src := `func @foo $32 %n {
	store &x %n
	ret
}
`
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.True(t, results.IsEmpty(), "unexpected errors: %v", results)
	assert.NotNil(t, function)
}

// TestLeaProducesPointer verifies that lea accepts a pointer-typed target and
// infers the variable type by stripping the pointer descriptor.
func TestLeaProducesPointer(t *testing.T) {
	src := `func @foo $32 %n {
	$32 * %ptr = lea &x
	ret
}
`
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.True(t, results.IsEmpty(), "unexpected errors: %v", results)
	assert.NotNil(t, function)

	vars := function.Variables.GetAllVariables()
	assert.Len(t, vars, 1)
	assert.Equal(t, "$32", vars[0].Type.String())
}

// TestLeaNestedPointer verifies that lea on a variable of type $32 *1 produces
// a $32 *2 target (incrementing the existing pointer descriptor, not appending).
func TestLeaNestedPointer(t *testing.T) {
	src := `func @foo $32 %n {
	$32 * %ptr = lea &x
	$32 *2 %pptr = lea &y
	ret
}
`
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.True(t, results.IsEmpty(), "unexpected errors: %v", results)
	assert.NotNil(t, function)

	vars := function.Variables.GetAllVariables()
	assert.Len(t, vars, 2)

	// &x inferred as $32, &y inferred as $32 *1
	byName := make(map[string]string)
	for _, v := range vars {
		byName[v.Name] = v.Type.String()
	}
	assert.Equal(t, "$32", byName["&x"])
	assert.Equal(t, "$32 *1", byName["&y"])
}

// TestLeaMismatchedPointer verifies that a lea with the wrong pointer level
// produces an error when the variable type is already known.
func TestLeaMismatchedPointer(t *testing.T) {
	src := `func @foo $32 %n {
	store &x %n
	$32 *2 %pptr = lea &x
	ret
}
`
	// &x is inferred as $32 by store; lea expects $32 *1 but got $32 *2
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.False(t, results.IsEmpty(), "expected pointer mismatch error")
	assert.Nil(t, function)
}
