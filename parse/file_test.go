package parse_test

import (
	"testing"
	"usm/lex"
	"usm/parse"
	"usm/source"

	"github.com/stretchr/testify/assert"
)

// The purpose of the test is to verify the structure of "simple" source file.
func TestSingleFunction(t *testing.T) {
	src := `def $i32 @add $i32 %x $i32 %y {
	%res = add %x %y
	ret %res
}`
	v := source.NewSourceView(src)
	srcView, _ := v.Detach()

	expected := parse.FileNode{
		Functions: []parse.FunctionNode{
			parse.FunctionNode{
				UnmanagedSourceView: srcView,
				Signature: parse.SignatureNode{
					UnmanagedSourceView: srcView.Subview(4, 29),
					Identifier:          srcView.Subview(9, 13),
					Parameters: []parse.ParameterNode{
						parse.ParameterNode{
							Type:     parse.TypeNode{srcView.Subview(14, 18)},
							Register: parse.RegisterNode{srcView.Subview(19, 21)},
						},
						parse.ParameterNode{
							Type:     parse.TypeNode{srcView.Subview(22, 26)},
							Register: parse.RegisterNode{srcView.Subview(27, 29)},
						},
					},
					Returns: []parse.TypeNode{
						parse.TypeNode{srcView.Subview(4, 8)},
					},
				},
				Block: parse.BlockNode{
					UnmanagedSourceView: srcView.Subview(30, 61),
					Instructions: []parse.InstructionNode{
						parse.InstructionNode{
							Operator: srcView.Subview(40, 43),
							Arguments: []parse.ArgumentNode{
								parse.ArgumentNode{srcView.Subview(44, 46)},
								parse.ArgumentNode{srcView.Subview(47, 49)},
							},
							Targets: []parse.RegisterNode{
								parse.RegisterNode{srcView.Subview(33, 37)},
							},
						},
						parse.InstructionNode{
							Operator: srcView.Subview(51, 54),
							Arguments: []parse.ArgumentNode{
								parse.ArgumentNode{srcView.Subview(55, 59)},
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
	file, perr := parse.FileParser{}.Parse(&tknsView)

	assert.Nil(t, perr)
	assert.Equal(t, expected, file)
}

func TestFileParserTwoFunctionsNoExtraSeparator(t *testing.T) {
	src := `def @first {
}
def @second {
}`

	expected := `def @first {
}

def @second {
}
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
