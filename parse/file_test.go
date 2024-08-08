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
	src := `function $32 @add $32 %x $32 %y =
	%res = add %x %y
	ret %res`
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
				Instructions: []parse.InstructionNode{
					{
						Operator: srcView.Subview(42, 45),
						Arguments: []parse.ArgumentNode{
							{srcView.Subview(46, 48)},
							{srcView.Subview(49, 51)},
						},
						Targets: []parse.RegisterNode{
							{srcView.Subview(35, 39)},
						},
					},
					{
						Operator: srcView.Subview(53, 56),
						Arguments: []parse.ArgumentNode{
							{srcView.Subview(57, 61)},
						},
					},
				},
			},
		},
	}

	tkns, err := lex.NewTokenizer().Tokenize(v)
	assert.NoError(t, err)

	tknsView := parse.NewTokenView(tkns)
	file, perr := parse.FileParser{}.Parse(&tknsView)

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
