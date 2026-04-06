package usmssa_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	usmmanagers "alon.kr/x/usm/usm/managers"
	usmssa "alon.kr/x/usm/usm/ssa"
	"github.com/stretchr/testify/assert"
)

const (
	inputFuncName    = "@input"
	expectedFuncName = "@expected"
)

func generateFileInfo(t *testing.T, source string) *gen.FileInfo {
	t.Helper()

	srcView := core.NewSourceView(source)
	tkns, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	tknView := parse.NewTokenView(tkns)
	fileNode, result := parse.NewFileParser().Parse(&tknView)
	assert.Nil(t, result)

	ctx := usmmanagers.NewGenerationContext()
	generator := gen.NewFileGenerator()
	info, results := generator.Generate(ctx, srcView.Ctx(), fileNode)
	assert.True(t, results.IsEmpty(), "Failed to generate file info")

	return info
}

func TestSsaDestruction(t *testing.T) {
	basePath := filepath.Join("testdata", "ssa_destruction")
	entries, err := os.ReadDir(basePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, entries)

	for _, entry := range entries {
		testName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		t.Run(testName, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(basePath, entry.Name()))
			assert.NoError(t, err)

			file := generateFileInfo(t, string(content))
			inputFunc := file.GetFunction(inputFuncName)
			assert.NotNil(t, inputFunc)
			expectedFunc := file.GetFunction(expectedFuncName)
			assert.NotNil(t, expectedFunc)

			results := usmssa.FunctionOutOfSsaForm(inputFunc)
			assert.True(t, results.IsEmpty(), "SSA destruction returned errors")

			inputFunc.Name = expectedFuncName
			assert.Equal(t, expectedFunc.String(), inputFunc.String())
		})
	}
}
