package gen_test

import (
	"math/big"
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func TestImmediateValueArgument(t *testing.T) {
	src := core.NewSourceView("$32 #1337")
	unmanaged := src.Unmanaged()
	node := parse.ImmediateNode{
		Type: parse.TypeNode{
			Identifier: unmanaged.Subview(0, 3),
		},
		Value: parse.ImmediateFinalValueNode{
			UnmanagedSourceView: unmanaged.Subview(4, 9),
		},
	}

	intType := &gen.NamedTypeInfo{Name: "$32", Size: 4}
	types := TypeMap{intType.Name: intType}

	ctx := gen.InstructionGenerationContext{
		FunctionGenerationContext: &gen.FunctionGenerationContext{
			FileGenerationContext: &gen.FileGenerationContext{
				GenerationContext: &testGenerationContext,
				SourceContext:     src.Ctx(),
				Types:             &types,
			},
			Registers: &RegisterMap{},
		},
		InstructionInfo: gen.NewEmptyInstructionInfo(&unmanaged),
	}

	generator := gen.NewImmediateArgumentGenerator()
	argument, results := generator.Generate(&ctx, node)

	assert.True(t, results.IsEmpty())
	assert.Equal(t, intType, argument.GetType().Base)

	immediateArgument, ok := argument.(*gen.ImmediateInfo)
	assert.True(t, ok)
	assert.Zero(t, immediateArgument.Value.Cmp(big.NewInt(1337)))
}
