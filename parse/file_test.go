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
	src := `function $32 @add $32 %x $32 %y = {
	%res = add %x %y
	ret %res
}`
	v := source.NewSourceView(src)
	srcView, _ := v.Detach()

	expected := parse.FileNode{
		Functions: []parse.FunctionNode{
			{
				UnmanagedSourceView: srcView,
				Declaration: parse.FunctionDeclarationNode{
					UnmanagedSourceView: srcView.Subview(9, 31),
					Identifier:          srcView.Subview(13, 17),
					Parameters: []parse.ParameterNode{
						{
							Type:     parse.TypeNode{srcView.Subview(18, 21)},
							Register: parse.RegisterNode{srcView.Subview(22, 24)},
						},
						{
							Type:     parse.TypeNode{srcView.Subview(25, 28)},
							Register: parse.RegisterNode{srcView.Subview(29, 31)},
						},
					},
					Returns: []parse.TypeNode{
						{srcView.Subview(9, 12)},
					},
				},
				Instructions: parse.BlockNode[parse.InstructionNode]{
					UnmanagedSourceView: srcView.Subview(34, 75),
					Nodes: []parse.InstructionNode{
						{
							Operator: srcView.Subview(44, 47),
							Arguments: []parse.ArgumentNode{
								parse.RegisterNode{srcView.Subview(48, 50)},
								parse.RegisterNode{srcView.Subview(51, 53)},
							},
							Targets: []parse.RegisterNode{
								{srcView.Subview(37, 41)},
							},
						},
						{
							Operator: srcView.Subview(55, 58),
							Arguments: []parse.ArgumentNode{
								parse.RegisterNode{srcView.Subview(59, 63)},
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
	src := `function @first =
		ret
function @second =
ret`

	expected := `function @first =
	ret

function @second =
	ret
`
	testFormattedFile(t, src, expected)
}

func TestFileWithLabels(t *testing.T) {
	src := `function $32 @fib $i32 %n =
    jle %n $32 #1 .return
    %n = dec %n
    %a = call @fib %n
    %n = dec %n
    %b = call @fib %n
    %n = add %a %b
.return
    ret %n`

	expected := `function $32 @fib $i32 %n =
	jle %n $32 #1 .return
	%n = dec %n
	%a = call @fib %n
	%n = dec %n
	%b = call @fib %n
	%n = add %a %b
	.return ret %n
`
	testFormattedFile(t, src, expected)
}

// MARK: Helpers

func testFormattedFile(t *testing.T, src string, expected string) {
	v := source.NewSourceView(src)
	_, ctx := v.Detach()
	tkns, err := lex.NewTokenizer().Tokenize(v)
	assert.NoError(t, err)

	tknsView := parse.NewTokenView(tkns)
	file, perr := parse.FileParser{}.Parse(&tknsView)
	assert.Nil(t, perr)

	got := file.String(ctx)
	assert.Equal(t, expected, got)
}
