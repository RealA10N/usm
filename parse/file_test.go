package parse_test

import (
	"testing"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"

	"github.com/stretchr/testify/assert"
)

// The purpose of the test is to verify the structure of "simple" source file.
func TestSingleFunction(t *testing.T) {
	src := `func $32 @add $32 %x $32 %y {
	%res = add %x %y
	ret %res
}`
	v := source.NewSourceView(src)
	srcView, _ := v.Detach()

	expected := parse.FileNode{
		Functions: []parse.FunctionNode{
			{
				UnmanagedSourceView: srcView,
				Signature: parse.FunctionSignatureNode{
					UnmanagedSourceView: srcView.Subview(5, 27),
					Identifier:          srcView.Subview(9, 13),
					Parameters: []parse.ParameterNode{
						{
							Type:     parse.TypeNode{Identifier: srcView.Subview(14, 17)},
							Register: parse.RegisterNode{srcView.Subview(18, 20)},
						},
						{
							Type:     parse.TypeNode{Identifier: srcView.Subview(21, 24)},
							Register: parse.RegisterNode{srcView.Subview(25, 27)},
						},
					},
					Returns: []parse.TypeNode{
						{Identifier: srcView.Subview(5, 8)},
					},
				},
				Instructions: &parse.BlockNode[parse.InstructionNode]{
					UnmanagedSourceView: srcView.Subview(28, 69),
					Nodes: []parse.InstructionNode{
						{
							Operator: srcView.Subview(38, 41),
							Arguments: []parse.ArgumentNode{
								parse.RegisterNode{srcView.Subview(42, 44)},
								parse.RegisterNode{srcView.Subview(45, 47)},
							},
							Targets: []parse.RegisterNode{
								{srcView.Subview(31, 35)},
							},
						},
						{
							Operator: srcView.Subview(49, 52),
							Arguments: []parse.ArgumentNode{
								parse.RegisterNode{srcView.Subview(53, 57)},
							},
						},
					},
				},
			},
		},
	}

	tkns, err := lex.NewTokenizer().Tokenize(v)
	assert.NoError(t, err)

	tknsView := parse.NewTokenView(tkns)
	file, perr := parse.NewFileParser().Parse(&tknsView)

	assert.Nil(t, perr)
	assert.Equal(t, expected, file)
}

func TestFileParserTwoFunctionsNoExtraSeparator(t *testing.T) {
	src := `func @first { ret }
func @second   {
ret
	}`

	expected := `func @first {
	ret
}

func @second {
	ret
}
`
	testExpectedFileFormat(t, src, expected)
}

func TestFileWithLabels(t *testing.T) {
	src := `
func $32 @fib $i32 %n {

	jle %n $32 #1 .return
    %n = dec %n
    %a = call @fib %n
    %n = dec %n
		%b = call @fib %n
    %n = add %a %b

.return ret %n

	}
`

	expected := `func $32 @fib $i32 %n {
	jle %n $32 #1 .return
	%n = dec %n
	%a = call @fib %n
	%n = dec %n
	%b = call @fib %n
	%n = add %a %b
.return
	ret %n
}
`
	testExpectedFileFormat(t, src, expected)
}

func TestFileWithTypes(t *testing.T) {
	src := `func @foo {
	ret
}

type $bar {
	$32 *
	.a $32 ^100
}

func @baz {
	ret
}

type $qux {
	$8 ^1337
}
`
	testExpectedFileFormat(t, src, src)
}

// MARK: Helpers

func testExpectedFileFormat(t *testing.T, src string, expected string) {
	t.Helper()

	v := source.NewSourceView(src)
	tkns, err := lex.NewTokenizer().Tokenize(v)
	assert.NoError(t, err)

	tknsView := parse.NewTokenView(tkns)
	file, perr := parse.NewFileParser().Parse(&tknsView)
	assert.Nil(t, perr)

	strCtx := parse.StringContext{SourceContext: v.Ctx()}
	assert.Equal(t, expected, file.String(&strCtx))
}
