package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type NamedTypeGenerator struct {
	ReferencedTypeGenerator FileContextGenerator[parse.TypeNode, ReferencedTypeInfo]
}

func NewNamedTypeGenerator() FileContextGenerator[parse.TypeDeclarationNode, *NamedTypeInfo] {
	return FileContextGenerator[parse.TypeDeclarationNode, *NamedTypeInfo](
		&NamedTypeGenerator{
			ReferencedTypeGenerator: NewReferencedTypeGenerator(),
		},
	)
}

func (g *NamedTypeGenerator) calculateTypeSize(
	ctx *FileGenerationContext,
	node parse.TypeNode,
	typeInfo ReferencedTypeInfo,
) (core.UsmUint, core.ResultList) {
	size := core.UsmUint(typeInfo.Base.Size)

	for _, descriptor := range typeInfo.Descriptors {
		switch descriptor.Type {
		case PointerTypeDescriptor:
			size = ctx.PointerSize
		case RepeatTypeDescriptor:
			var ok bool
			size, ok = core.Mul(size, descriptor.Amount)
			if !ok {
				v := node.View()
				return 0, list.FromSingle(core.Result{{
					Type:     core.ErrorResult,
					Message:  "Type size overflow",
					Location: &v,
				}})
			}
		default:
			// notest
			v := node.View()
			return 0, list.FromSingle(core.Result{{
				Type:     core.InternalErrorResult,
				Message:  "Unknown type descriptor",
				Location: &v,
			}})
		}
	}

	return size, core.ResultList{}
}

func (g *NamedTypeGenerator) Generate(
	ctx *FileGenerationContext,
	node parse.TypeDeclarationNode,
) (*NamedTypeInfo, core.ResultList) {
	identifier := viewToSourceString(ctx, node.Identifier)
	declaration := node.View()

	typeInfo := ctx.Types.GetType(identifier)
	if typeInfo != nil {
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Type already defined",
				Location: &declaration,
			},
			{
				Type:     core.HintResult,
				Message:  "Previous definition here",
				Location: typeInfo.Declaration,
			},
		})
	}

	if len(node.Fields.Nodes) != 1 {
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Type declaration currently only supports a single field",
				Location: &declaration,
			},
		})
	}

	referencedTypeNode := node.Fields.Nodes[0].Type
	referencedTypeInfo, results := g.ReferencedTypeGenerator.Generate(ctx, referencedTypeNode)
	if !results.IsEmpty() {
		return nil, results
	}

	size, results := g.calculateTypeSize(ctx, referencedTypeNode, referencedTypeInfo)
	if !results.IsEmpty() {
		return nil, results
	}

	typeInfo = &NamedTypeInfo{
		Name:        identifier,
		Size:        size,
		Declaration: &declaration,
	}

	result := ctx.Types.NewType(typeInfo)
	if result != nil {
		return nil, list.FromSingle(result)
	}

	return typeInfo, core.ResultList{}
}
