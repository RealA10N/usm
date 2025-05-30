package gen_test

import (
	"math/big"
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func TestTypeAliasDeclaration(t *testing.T) {
	typeManager := make(TypeMap)
	typeManager.newBuiltinType("$32", big.NewInt(32))

	view := core.NewSourceView("type $myType $32")
	unmanaged := view.Unmanaged()

	typeDeclarationNode := parse.TypeDeclarationNode{
		UnmanagedSourceView: unmanaged,
		Identifier:          unmanaged.Subview(5, 12),
		Fields: parse.BlockNode[parse.TypeFieldNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{},
			Nodes: []parse.TypeFieldNode{
				{
					Type: parse.TypeNode{
						Identifier: unmanaged.Subview(13, 16),
					},
				},
			},
		},
	}

	ctx := gen.FileGenerationContext{
		GenerationContext: &testGenerationContext,
		SourceContext:     view.Ctx(),
		Types:             &typeManager,
	}

	generator := gen.NewNamedTypeGenerator()
	typeInfo, results := generator.Generate(&ctx, typeDeclarationNode)

	assert.True(t, results.IsEmpty())
	assert.NotNil(t, typeInfo)
	assert.Equal(t, "$myType", string(typeInfo.Name))
	assert.Zero(t, typeInfo.Size.Cmp(big.NewInt(32)))
}

func TestPointerTypeDeclaration(t *testing.T) {
	typeManager := make(TypeMap)
	typeManager.newBuiltinType("$8", big.NewInt(8))

	view := core.NewSourceView("type $myType $8 *")
	unmanaged := view.Unmanaged()

	typeDeclarationNode := parse.TypeDeclarationNode{
		UnmanagedSourceView: unmanaged,
		Identifier:          unmanaged.Subview(5, 12),
		Fields: parse.BlockNode[parse.TypeFieldNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{},
			Nodes: []parse.TypeFieldNode{
				{
					Type: parse.TypeNode{
						Identifier: unmanaged.Subview(13, 15),
						Decorators: []parse.TypeDecoratorNode{
							{
								UnmanagedSourceView: unmanaged.Subview(16, 17),
								Type:                parse.PointerTypeDecorator,
							},
						},
					},
				},
			},
		},
	}

	ctx := gen.FileGenerationContext{
		GenerationContext: &testGenerationContext,
		SourceContext:     view.Ctx(),
		Types:             &typeManager,
	}

	generator := gen.NewNamedTypeGenerator()
	typeInfo, results := generator.Generate(&ctx, typeDeclarationNode)

	assert.True(t, results.IsEmpty())
	assert.NotNil(t, typeInfo)
	assert.Equal(t, "$myType", string(typeInfo.Name))
	assert.EqualValues(t, testGenerationContext.PointerSize, typeInfo.Size)
}

func TestRepeatTypeDeclaration(t *testing.T) {
	typeManager := make(TypeMap)
	typeManager.newBuiltinType("$8", big.NewInt(8))

	view := core.NewSourceView("type $myType $8 ^9")
	unmanaged := view.Unmanaged()

	typeDeclarationNode := parse.TypeDeclarationNode{
		UnmanagedSourceView: unmanaged,
		Identifier:          unmanaged.Subview(5, 12),
		Fields: parse.BlockNode[parse.TypeFieldNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{},
			Nodes: []parse.TypeFieldNode{
				{
					Type: parse.TypeNode{
						Identifier: unmanaged.Subview(13, 15),
						Decorators: []parse.TypeDecoratorNode{
							{
								UnmanagedSourceView: unmanaged.Subview(16, 18),
								Type:                parse.RepeatTypeDecorator,
							},
						},
					},
				},
			},
		},
	}

	ctx := gen.FileGenerationContext{
		GenerationContext: &testGenerationContext,
		SourceContext:     view.Ctx(),
		Types:             &typeManager,
	}

	generator := gen.NewNamedTypeGenerator()
	typeInfo, results := generator.Generate(&ctx, typeDeclarationNode)

	assert.True(t, results.IsEmpty())
	assert.NotNil(t, typeInfo)
	assert.Equal(t, "$myType", typeInfo.Name)
	assert.Zero(t, typeInfo.Size.Cmp(big.NewInt(8*9)))
}

func TestAlreadyDefinedTypeDeclaration(t *testing.T) {
	typeManager := make(TypeMap)
	typeManager.newBuiltinType("$32", big.NewInt(32))

	intTypeInfo := gen.NewNamedTypeInfo("$int", big.NewInt(32), nil)
	typeManager.NewType(intTypeInfo)

	view := core.NewSourceView("type $int $32")
	unmanaged := view.Unmanaged()

	node := parse.TypeDeclarationNode{
		UnmanagedSourceView: unmanaged,
		Identifier:          unmanaged.Subview(5, 9),
		Fields: parse.BlockNode[parse.TypeFieldNode]{
			Nodes: []parse.TypeFieldNode{
				{
					Type: parse.TypeNode{
						Identifier: unmanaged.Subview(10, 13),
						Decorators: []parse.TypeDecoratorNode{},
					},
				},
			},
		},
	}

	genCtx := gen.FileGenerationContext{
		GenerationContext: &testGenerationContext,
		SourceContext:     view.Ctx(),
		Types:             &typeManager,
	}

	generator := gen.NewNamedTypeGenerator()
	_, results := generator.Generate(&genCtx, node)

	assert.Len(t, results.ToSlice(), 1)
	result := results.Head.Value
	assert.Contains(t, result[0].Message, "already defined")
}

// TODO: add back this test when more complex types are supported!
// func TestVoidTypeDeclaration(t *testing.T) {
// 	typeManager := make(TypeMap)
// 	typeDeclarationNode := parse.TypeDeclarationNode{
// 		Fields: parse.BlockNode[parse.TypeFieldNode]{
// 			UnmanagedSourceView: core.UnmanagedSourceView{},
// 			Nodes:               []parse.TypeFieldNode{},
// 		},
// 	}

// 	genCtx := gen.GenerationContext[gen.BaseInstruction]{
// 		Types: &typeManager,
// 	}

// 	generator := gen.NewNamedTypeGenerator[gen.BaseInstruction]()
// 	typeInfo, results := generator.Generate(&genCtx, typeDeclarationNode)

// 	assert.True(t, results.IsEmpty())
// 	assert.NotNil(t, typeInfo)
// 	assert.EqualValues(t, 0, typeInfo.Size)
// }

func TestRepeatTypeTooLarge(t *testing.T) {
	typeManager := make(TypeMap)
	typeManager.newBuiltinType("$32", big.NewInt(32))

	v := core.NewSourceView("type $tooLarge { $32 ^1_000_000_000 ^1_000_000_000 }")
	unmanaged := v.Unmanaged()

	node := parse.TypeDeclarationNode{
		UnmanagedSourceView: unmanaged,
		Identifier:          unmanaged.Subview(5, 14),
		Fields: parse.BlockNode[parse.TypeFieldNode]{
			UnmanagedSourceView: unmanaged.Subview(15, 44),
			Nodes: []parse.TypeFieldNode{
				{
					Type: parse.TypeNode{
						Identifier: unmanaged.Subview(17, 20),
						Decorators: []parse.TypeDecoratorNode{
							{
								UnmanagedSourceView: unmanaged.Subview(21, 31),
								Type:                parse.RepeatTypeDecorator,
							},
							{
								UnmanagedSourceView: unmanaged.Subview(32, 42),
								Type:                parse.RepeatTypeDecorator,
							},
						},
					},
				},
			},
		},
	}

	ctx := &gen.FileGenerationContext{
		GenerationContext: &testGenerationContext,
		SourceContext:     v.Ctx(),
		Types:             &typeManager,
	}

	generator := gen.NewNamedTypeGenerator()
	_, results := generator.Generate(ctx, node)
	assert.False(t, results.IsEmpty())
}
