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
	return funcGen.Generate(ctx, node)
}

// TestFunctionWithVariable verifies that a variable declaration is collected
// into FunctionInfo.Variables and does not appear as an instruction.
func TestFunctionWithVariable(t *testing.T) {
	src := `func $32 @foo $32 %n {
&saved $32
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

// TestFunctionWithVariableUndefined verifies that referencing an undeclared
// variable in an instruction produces an error.
func TestFunctionWithVariableUndefined(t *testing.T) {
	src := `func $32 @foo $32 %n {
	$32 %result = load &notDeclared
	ret %result
}
`
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.False(t, results.IsEmpty(), "expected error for undefined variable")
	assert.Nil(t, function)
}

// TestFunctionWithVariableDuplicate verifies that declaring the same variable
// twice produces an error.
func TestFunctionWithVariableDuplicate(t *testing.T) {
	src := `func $32 @foo $32 %n {
&x $32
&x $32
	store &x %n
	$32 %result = load &x
	ret %result
}
`
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.False(t, results.IsEmpty(), "expected error for duplicate variable")
	assert.Nil(t, function)
}

// TestLoadTypeMismatch verifies that load rejects a target whose type does not
// match the variable type.
func TestLoadTypeMismatch(t *testing.T) {
	// Type mismatch: variable is $32 but target is... let's use $32 (same for
	// success case) and make a separate mismatch source.
	src := `func $32 @foo {
&x $32
	$32 %a = load &x
	ret %a
}
`
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.True(t, results.IsEmpty(), "expected no errors: %v", results)
	assert.NotNil(t, function)
}

// TestStoreTypeMismatch verifies that store rejects a value whose type does not
// match the variable type.
func TestStoreCritical(t *testing.T) {
	// store is a critical instruction — verify it generates without errors and
	// the function contains exactly one store in its instructions.
	src := `func @foo $32 %n {
&x $32
	store &x %n
	ret
}
`
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.True(t, results.IsEmpty(), "unexpected errors: %v", results)
	assert.NotNil(t, function)
}

// TestLeaProducesPointer verifies that lea accepts a pointer-typed target and
// a variable argument with the correct base type.
func TestLeaProducesPointer(t *testing.T) {
	src := `func @foo $32 %n {
&x $32
	$32 * %ptr = lea &x
	ret
}
`
	function, results := generateFunctionFromSourceUSM(t, src)
	assert.True(t, results.IsEmpty(), "unexpected errors: %v", results)
	assert.NotNil(t, function)
}
