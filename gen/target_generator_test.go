package gen_test

import (
	"math/big"
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
			TokenNode: parse.TokenNode{
				UnmanagedSourceView: unmanaged.Subview(4, 6),
			},
		},
		Type: &parse.TypeNode{
			Identifier: unmanaged.Subview(0, 3),
		},
	}

	intType := gen.NewNamedTypeInfo("$32", big.NewInt(32), nil)
	types := TypeMap{intType.Name: intType}

	intTypeRef := gen.ReferencedTypeInfo{
		Base: intType,
	}

	registers := RegisterMap{
		"%a": &gen.RegisterInfo{Name: "%a", Type: intTypeRef},
	}

	ctx := &gen.FunctionGenerationContext{
		FileGenerationContext: &gen.FileGenerationContext{
			GenerationContext: &testGenerationContext,
			SourceContext:     src.Ctx(),
			Types:             &types,
		},
		Registers: &registers,
	}

	generator := gen.NewTargetGenerator()
	info, results := generator.Generate(ctx, node)
	assert.True(t, results.IsEmpty())
	assert.Equal(t, intType, info.Register.Type.Base)
}
