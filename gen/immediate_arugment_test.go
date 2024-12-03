package gen_test

import (
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

	intType := &gen.TypeInfo{Name: "$32", Size: 4}
	types := TypeMap{intType.Name: intType}

	ctx := gen.GenerationContext[Instruction]{
		SourceContext: src.Ctx(),
		Types:         &types,
		Registers:     &RegisterMap{},
		Instructions:  &InstructionMap{},
	}

	generator := gen.ImmediateArgumentGenerator[Instruction]{}
	_, results := generator.Generate(&ctx, node)
	assert.True(t, results.IsEmpty())
}
