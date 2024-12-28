package gen_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func TestUndefinedRegisterArgument(t *testing.T) {
	src := core.NewSourceView("%a")
	unmanaged := src.Unmanaged()
	node := parse.RegisterNode{TokenNode: parse.TokenNode{UnmanagedSourceView: unmanaged}}

	ctx := gen.FunctionGenerationContext{
		FileGenerationContext: &gen.FileGenerationContext{
			GenerationContext: &gen.GenerationContext{
				Instructions: &InstructionMap{},
				PointerSize:  8,
			},
			SourceContext: src.Ctx(),
			Types:         &TypeMap{},
		},
		Registers: &RegisterMap{},
	}

	generator := gen.NewArgumentGenerator()
	_, results := generator.Generate(&ctx, node)

	assert.EqualValues(t, 1, results.Len())
	result := results.Head.Value
	expectedResult := gen.UndefinedRegisterResult(node).Head.Value
	assert.EqualValues(t, expectedResult, result)
}
