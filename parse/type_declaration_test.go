package parse_test

import (
	"testing"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"
	"github.com/stretchr/testify/assert"
)

func TestTypeDeclarationVoid(t *testing.T) {
	src := "type $void { }"
	testExpectedTypeDeclaration(t, src, src)
}

func TestTypeDeclarationBool(t *testing.T) {
	src := "type $bool { $1 }"
	expected := "type $bool {\n\t$1\n}"
	testExpectedTypeDeclaration(t, src, expected)
}

func TestTypeDeclarationMultipleFieldsWithLabels(t *testing.T) {
	src := `type $strList {
	.str $8 *
	.next $strList *
	.prev $strList *
}`
	testExpectedTypeDeclaration(t, src, src)
}

func TestTypeDeclarationMultipleLabels(t *testing.T) {
	src := `type $struct {
	.a
			.b
.c

	$8 *2102 ^1337 * ^ * .d .e .f $1234
  }`

	expected := `type $struct {
	.a .b .c $8 *2102 ^1337 * ^ *
	.d .e .f $1234
}`
	testExpectedTypeDeclaration(t, src, expected)
}

func TestTypeDeclarationLongOneLine(t *testing.T) {
	src := `type $struct { .a .b .c $8 *2102 ^1337 * ^ * .d .e .f $1234 }`
	expected := `type $struct {
	.a .b .c $8 *2102 ^1337 * ^ *
	.d .e .f $1234
}`
	testExpectedTypeDeclaration(t, src, expected)
}

// MARK: Helpers

func testExpectedTypeDeclaration(t *testing.T, src, expected string) {
	t.Helper()

	srcView := source.NewSourceView(src)
	ctx := srcView.Ctx()
	tkns, err := lex.NewTokenizer().Tokenize(srcView)
	assert.NoError(t, err)

	tknView := parse.NewTokenView(tkns)
	p := parse.NewTypeDeclarationParser()
	typ, perr := p.Parse(&tknView)
	assert.Nil(t, perr)

	assert.Equal(t, expected, typ.String(&parse.StringContext{SourceContext: ctx}))
}
