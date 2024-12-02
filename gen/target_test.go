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

	intType := &gen.TypeInfo{Name: "$32", Size: 4}
	types := TypeMap{intType.Name: intType}

	registers := RegisterMap{
		"%a": &gen.RegisterInfo{Name: "%a", Type: intType},
	}

	ctx := gen.GenerationContext[Instruction]{
		ArchInfo:      gen.ArchInfo{PointerSize: 8},
		SourceContext: src.Ctx(),
		Types:         &types,
		Registers:     &registers,
		Instructions:  &InstructionMap{},
	}

	generator := gen.TargetGenerator[Instruction]{}
	info, results := generator.Generate(&ctx, node)
	assert.True(t, results.IsEmpty())
	assert.Equal(t, intType, info)
}
