package opt_test

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
	"github.com/stretchr/testify/assert"
)

const (
	InputFuncName    = "@input"
	ExpectedFuncName = "@expected"
)

// OptimizationTestFunc is a function that applies an optimization to a function.
type OptimizationTestFunc func(*gen.FunctionInfo) core.ResultList

// readSourceFile reads a source file and returns its contents.
func readSourceFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	assert.NoError(t, err, "Failed to read file: %s", path)

	return string(content)
}

// findTestPaths finds all test cases in the given directory.
// The directory structure is flat:
// testdata/optimization_name/testname.usm
// Each test file should contain exactly two functions: @input and @expected
func findTestPaths(t *testing.T, optimizationName string) []string {
	t.Helper()

	basePath := filepath.Join("testdata", optimizationName)
	entries, err := os.ReadDir(basePath)
	assert.NoError(t, err, "Failed to read directory: %s", basePath)

	var testPaths []string

	// Find all test files
	for _, entry := range entries {
		if entry.IsDir() {
			t.Fatalf("Unexpected nested directory found: %s", entry.Name())
		}

		testPath := filepath.Join(basePath, entry.Name())
		testPaths = append(testPaths, testPath)
	}

	return testPaths
}

// generateFileInfo parses and generates a file internal representation from the
// source code.
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

// extractTestFunctions extracts the @input and @expected functions from the
// file info.
func extractTestFunctions(
	t *testing.T,
	file *gen.FileInfo,
) (*gen.FunctionInfo, *gen.FunctionInfo) {
	t.Helper()
	inputFunc := file.GetFunction(InputFuncName)
	assert.NotNil(t, inputFunc, "Failed to extract @input function")

	expectedFunc := file.GetFunction(ExpectedFuncName)
	assert.NotNil(t, expectedFunc, "Failed to extract @expected function")

	return inputFunc, expectedFunc
}

// RunOptimizationTests runs all test cases for the given optimization.
func RunOptimizationTests(
	t *testing.T,
	optimizationName string,
	optimizationFunc OptimizationTestFunc,
) {
	testPaths := findTestPaths(t, optimizationName)
	assert.NotEmpty(t, testPaths, "No test cases found for optimization: %s", optimizationName)

	for _, testPath := range testPaths {
		testName := filepath.Base(testPath)
		testName = strings.TrimSuffix(testName, filepath.Ext(testName))

		t.Run(testName, func(t *testing.T) {
			source := readSourceFile(t, testPath)

			// Extract input and expected functions
			file := generateFileInfo(t, source)
			inputFunc, expectedFunc := extractTestFunctions(t, file)

			// Apply optimization to the input function
			results := optimizationFunc(inputFunc)
			assert.True(t, results.IsEmpty(), "Optimization returned errors")

			// Compare the optimized function with the expected function
			inputFunc.Name = ExpectedFuncName
			assert.Equal(
				t,
				expectedFunc.String(),
				inputFunc.String(),
				"Function bodies should match after optimization",
			)
		})
	}
}
