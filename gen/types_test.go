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
	m.RegisterType(name, &gen.TypeInfo{Size: size})
}

func (m TypeManager) RegisterType(name string, typ *gen.TypeInfo) core.UsmError {
	if m[name] != nil {
		return core.GenericError{
			ErrorMessage:  "Type already defined",
			ErrorLocation: typ.Name,
		}
	}

	m[name] = typ
	return nil
}

func TestBasicTypeAliasDeclaration(t *testing.T) {
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
	assert.Equal(t, "$myType", string(typeInfo.Name.Raw(genCtx.SourceContext)))
	assert.EqualValues(t, 4, typeInfo.Size)
}

func TestBasicTypeAliasPointerDeclaration(t *testing.T) {
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
	assert.Equal(t, "$myType", string(typeInfo.Name.Raw(genCtx.SourceContext)))
	assert.EqualValues(t, 1337, typeInfo.Size)
}
