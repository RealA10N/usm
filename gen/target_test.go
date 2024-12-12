package gen_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func TestTargetRegisterAlreadyDefined(t *testing.T) {
	src := core.NewSourceView("$32 %a")
	unmanaged := src.Unmanaged()
	node := parse.TargetNode{
		Register: parse.RegisterNode{
			UnmanagedSourceView: unmanaged.Subview(4, 6),
		},
		Type: &parse.TypeNode{
			Identifier: unmanaged.Subview(0, 3),
		},
	}

	intType := &gen.NamedTypeInfo{Name: "$32", Size: 4}
	types := TypeMap{intType.Name: intType}

	intTypeRef := gen.ReferencedTypeInfo{
		Base: intType,
		Size: intType.Size,
	}

	registers := RegisterMap{
		"%a": &gen.RegisterInfo{Name: "%a", Type: intTypeRef},
	}

	ctx := gen.FunctionGenerationContext[Instruction]{
		FileGenerationContext: &gen.FileGenerationContext[Instruction]{
			GenerationContext: &testGenerationContext,
			SourceContext:     src.Ctx(),
			Types:             &types,
		},
		Registers: &registers,
	}

	generator := gen.NewTargetGenerator[Instruction]()
	info, results := generator.Generate(&ctx, node)
	assert.True(t, results.IsEmpty())
	assert.Equal(t, intType, info.Type.Base)
}
