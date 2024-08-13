package parse_test

import (
	"testing"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"
	"github.com/stretchr/testify/assert"
)

func TestImmediateValue(t *testing.T) {
	src := "$32 #1337"
	testExpectedImmediate(t, src, src)
}

func TestImmediateBlockOneLine(t *testing.T) {
	src := "$custom { .a #1337 .b #2902 }"
	expected := `$custom {
	.a #1337
	.b #2902
}`

	testExpectedImmediate(t, src, expected)
}

func TestImmediateBlockNested(t *testing.T) {
	src := "$outer { .a #1337 .b { .c #2902 #1234 } }"
	expected := `$outer {
	.a #1337
	.b {
		.c #2902
		#1234
	}
}`

	testExpectedImmediate(t, src, expected)
}

// MARK: Helpers

func testExpectedImmediate(t *testing.T, src, expected string) {
	t.Helper()

	srcView := source.NewSourceView(src)
	tkns, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	v := parse.NewTokenView(tkns)
	immediate, perr := parse.NewImmediateParser().Parse(&v)
	assert.Nil(t, perr)

	strCtx := parse.StringContext{SourceContext: srcView.Ctx()}
	assert.Equal(t, expected, immediate.String(&strCtx))
}
