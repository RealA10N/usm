package gen_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

type TypeManager map[string]*gen.TypeInfo

func (m TypeManager) GetType(name string) *gen.TypeInfo {
	return m[name]
}

func (m TypeManager) registerBuiltinType(name string, size core.UsmUint) {
	m.NewType(name, &gen.TypeInfo{Size: size})
}

func (m TypeManager) NewType(name string, typ *gen.TypeInfo) core.Result {
	if m[name] != nil {
		return &core.GenericResult{
			Type:     core.ErrorResult,
			Message:  "Type already defined",
			Location: &typ.Declaration,
		}
	}

	m[name] = typ
	return nil
}

func TestTypeAliasDeclaration(t *testing.T) {
	typeManager := make(TypeManager)
	typeManager.registerBuiltinType("$32", 4)

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

	genCtx := gen.GenerationContext{
		SourceContext: view.Ctx(),
		Types:         typeManager,
	}

	typeInfo, err := gen.TypeInfoFromTypeDeclaration(&genCtx, typeDeclarationNode)
	assert.Nil(t, err)
	assert.NotNil(t, typeInfo)
	assert.Equal(t, "$myType", string(typeInfo.Name))
	assert.EqualValues(t, 4, typeInfo.Size)
}

func TestPointerTypeDeclaration(t *testing.T) {
	typeManager := make(TypeManager)
	typeManager.registerBuiltinType("$64", 8)

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

	genCtx := gen.GenerationContext{
		SourceContext: view.Ctx(),
		ArchInfo:      gen.ArchInfo{PointerSize: 1337},
		Types:         typeManager,
	}

	typeInfo, err := gen.TypeInfoFromTypeDeclaration(&genCtx, typeDeclarationNode)
	assert.Nil(t, err)
	assert.NotNil(t, typeInfo)
	assert.Equal(t, "$myType", string(typeInfo.Name))
	assert.EqualValues(t, 1337, typeInfo.Size)
}

func TestRepeatTypeDeclaration(t *testing.T) {
	typeManager := make(TypeManager)
	typeManager.registerBuiltinType("$8", 1)

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

	genCtx := gen.GenerationContext{
		SourceContext: view.Ctx(),
		Types:         typeManager,
	}

	typeInfo, err := gen.TypeInfoFromTypeDeclaration(&genCtx, typeDeclarationNode)
	assert.Nil(t, err)
	assert.NotNil(t, typeInfo)
	assert.Equal(t, "$myType", typeInfo.Name)
	assert.EqualValues(t, 9, typeInfo.Size)
}

func TestVoidTypeDeclaration(t *testing.T) {
	typeManager := make(TypeManager)
	typeDeclarationNode := parse.TypeDeclarationNode{
		Fields: parse.BlockNode[parse.TypeFieldNode]{
			UnmanagedSourceView: core.UnmanagedSourceView{},
			Nodes:               []parse.TypeFieldNode{},
		},
	}

	genCtx := gen.GenerationContext{
		Types: typeManager,
	}

	typeInfo, err := gen.TypeInfoFromTypeDeclaration(&genCtx, typeDeclarationNode)
	assert.Nil(t, err)
	assert.NotNil(t, typeInfo)
	assert.EqualValues(t, 0, typeInfo.Size)
}
