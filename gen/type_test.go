package gen_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func TestTypeAliasDeclaration(t *testing.T) {
	typeManager := make(TypeMap)
	typeManager.newBuiltinType("$32", 4)

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

	genCtx := gen.GenerationContext[gen.BaseInstruction]{
		SourceContext: view.Ctx(),
		Types:         &typeManager,
	}

	typeInfo, err := gen.TypeInfoFromTypeDeclaration(&genCtx, typeDeclarationNode)
	assert.Nil(t, err)
	assert.NotNil(t, typeInfo)
	assert.Equal(t, "$myType", string(typeInfo.Name))
	assert.EqualValues(t, 4, typeInfo.Size)
}

func TestPointerTypeDeclaration(t *testing.T) {
	typeManager := make(TypeMap)
	typeManager.newBuiltinType("$64", 8)

	view := core.NewSourceView("type $myType $64 *")
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
						Decorators: []parse.TypeDecoratorNode{
							{
								UnmanagedSourceView: unmanaged.Subview(17, 18),
								Type:                parse.PointerTypeDecorator,
							},
						},
					},
				},
			},
		},
	}

	genCtx := gen.GenerationContext[gen.BaseInstruction]{
		SourceContext: view.Ctx(),
		ArchInfo:      gen.ArchInfo{PointerSize: 1337},
		Types:         &typeManager,
	}

	typeInfo, err := gen.TypeInfoFromTypeDeclaration(&genCtx, typeDeclarationNode)
	assert.Nil(t, err)
	assert.NotNil(t, typeInfo)
	assert.Equal(t, "$myType", string(typeInfo.Name))
	assert.EqualValues(t, 1337, typeInfo.Size)
}

func TestRepeatTypeDeclaration(t *testing.T) {
	typeManager := make(TypeMap)
	typeManager.newBuiltinType("$8", 1)

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

	genCtx := gen.GenerationContext[gen.BaseInstruction]{
		SourceContext: view.Ctx(),
		Types:         &typeManager,
	}

	typeInfo, err := gen.TypeInfoFromTypeDeclaration(&genCtx, typeDeclarationNode)
	assert.Nil(t, err)
	assert.NotNil(t, typeInfo)
	assert.Equal(t, "$myType", typeInfo.Name)
	assert.EqualValues(t, 9, typeInfo.Size)
}

func TestVoidTypeDeclaration(t *testing.T) {
	typeManager := make(TypeMap)
	typeDeclarationNode := parse.TypeDeclarationNode{
		Fields: parse.BlockNode[parse.TypeFieldNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{},
			Nodes:               []parse.TypeFieldNode{},
		},
	}

	genCtx := gen.GenerationContext[gen.BaseInstruction]{
		Types: &typeManager,
	}

	typeInfo, err := gen.TypeInfoFromTypeDeclaration(&genCtx, typeDeclarationNode)
	assert.Nil(t, err)
	assert.NotNil(t, typeInfo)
	assert.EqualValues(t, 0, typeInfo.Size)
}
