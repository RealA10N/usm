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
	usm64managers "alon.kr/x/usm/usm64/managers"
	"github.com/stretchr/testify/assert"
)

// OptimizationTestCase represents a test case for an optimization.
type OptimizationTestCase struct {
	Name         string
	InputPath    string
	ExpectedPath string
}

// OptimizationTestFunc is a function that applies an optimization to a function.
type OptimizationTestFunc func(*gen.FunctionInfo) core.ResultList

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

// readSourceFile reads a source file and returns its contents.
func readSourceFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	assert.NoError(t, err, "Failed to read file: %s", path)

	return string(content)
}

// findTestCases finds all test cases in the given directory.
// The directory structure is flat:
// testdata/optimization_name/testname_input.usm
// testdata/optimization_name/testname_expected.usm
func findTestCases(t *testing.T, optimizationName string) []OptimizationTestCase {
	t.Helper()

	basePath := filepath.Join("testdata", optimizationName)
	entries, err := os.ReadDir(basePath)
	assert.NoError(t, err, "Failed to read directory: %s", basePath)

	// Map to store test cases by name
	testCaseMap := make(map[string]OptimizationTestCase)

	// Find all input and expected files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		// Parse the file name to get the test name and type (input or expected)
		if strings.HasSuffix(fileName, "_input.usm") {
			testName := strings.TrimSuffix(fileName, "_input.usm")
			inputPath := filepath.Join(basePath, fileName)

			// Create or update the test case
			testCase, exists := testCaseMap[testName]
			if !exists {
				testCase = OptimizationTestCase{Name: testName}
			}
			testCase.InputPath = inputPath
			testCaseMap[testName] = testCase
		} else if strings.HasSuffix(fileName, "_expected.usm") {
			testName := strings.TrimSuffix(fileName, "_expected.usm")
			expectedPath := filepath.Join(basePath, fileName)

			// Create or update the test case
			testCase, exists := testCaseMap[testName]
			if !exists {
				testCase = OptimizationTestCase{Name: testName}
			}
			testCase.ExpectedPath = expectedPath
			testCaseMap[testName] = testCase
		}
	}

	// Convert the map to a slice of test cases
	var testCases []OptimizationTestCase
	for _, testCase := range testCaseMap {
		// Ensure each test case has both input and expected files
		if testCase.InputPath == "" {
			t.Errorf("Test case %q is missing input file", testCase.Name)
		} else if testCase.ExpectedPath == "" {
			t.Errorf("Test case %q is missing expected file", testCase.Name)
		} else {
			testCases = append(testCases, testCase)
		}
	}

	return testCases
}

// RunOptimizationTests runs all test cases for the given optimization.
func RunOptimizationTests(t *testing.T, optimizationName string, optimizationFunc OptimizationTestFunc) {
	testCases := findTestCases(t, optimizationName)
	assert.NotEmpty(t, testCases, "No test cases found for optimization: %s", optimizationName)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Read input and expected source
			inputSource := readSourceFile(t, tc.InputPath)
			expectedSource := readSourceFile(t, tc.ExpectedPath)

			// Generate function from input source
			function, results := generateFunctionFromSource(t, inputSource)
			assert.True(t, results.IsEmpty(), "Failed to generate function from input source")

			// Apply optimization
			results = optimizationFunc(function)
			assert.True(t, results.IsEmpty(), "Optimization returned errors")

			// Compare the optimized function with the expected function
			assert.Equal(t, expectedSource, function.String())
		})
	}
}
