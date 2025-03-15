package opt_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/opt"
	"alon.kr/x/usm/parse"
	usm64managers "alon.kr/x/usm/usm64/managers"
	"github.com/stretchr/testify/assert"
)

// generateFunctionFromSource is a test utility function that parses and generates
// a function from the given source code using the usm64 ISA.
func generateFunctionFromSource(
	t *testing.T,
	source string,
) (*gen.FunctionInfo, core.ResultList) {
	t.Helper()

	ctx := usm64managers.NewGenerationContext()
	src := core.NewSourceView(source)
	fileCtx := gen.CreateFileContext(ctx, src.Ctx())

	tkns, err := lex.NewTokenizer().Tokenize(src)
	assert.NoError(t, err)

	tknView := parse.NewTokenView(tkns)
	node, result := parse.NewFunctionParser().Parse(&tknView)
	assert.Nil(t, result)

	generator := gen.NewFunctionGenerator()
	return generator.Generate(fileCtx, node)
}

func TestBasicDeadCodeElimination(t *testing.T) {
	src := `func @main {
	$64 %a = $64 #0
	$64 %b = $64 #1
	PUT %a
	TERM
}`

	function, results := generateFunctionFromSource(t, src)
	assert.True(t, results.IsEmpty())

	results = opt.DeadCodeElimination(function)
	assert.True(t, results.IsEmpty())

	assert.Len(t, function.EntryBlock.Instructions, 3)
	assert.Equal(t, "$64 %a = $64 #0", function.EntryBlock.Instructions[0].String())
	assert.Equal(t, "PUT %a", function.EntryBlock.Instructions[1].String())
	assert.Equal(t, "TERM", function.EntryBlock.Instructions[2].String())
}
