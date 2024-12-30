package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type ReferencedTypeGenerator struct {
	DescriptorGenerator FileContextGenerator[parse.TypeDecoratorNode, TypeDescriptorInfo]
}

func NewReferencedTypeGenerator() FileContextGenerator[parse.TypeNode, ReferencedTypeInfo] {
	return FileContextGenerator[parse.TypeNode, ReferencedTypeInfo](
		&ReferencedTypeGenerator{
			DescriptorGenerator: NewDescriptorGenerator(),
		},
	)
}

func (g *ReferencedTypeGenerator) Generate(
	ctx *FileGenerationContext,
	node parse.TypeNode,
) (ReferencedTypeInfo, core.ResultList) {
	baseIdentifier := viewToSourceString(ctx, node.Identifier)
	baseType := ctx.Types.GetType(baseIdentifier)

	if baseType == nil {
		return ReferencedTypeInfo{}, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Undefined type",
				Location: &node.Identifier,
			},
			// TODO: add a hint in case of a typo hint?
		})
	}

	descriptors := make([]TypeDescriptorInfo, 0, len(node.Decorators))
	for _, descriptor := range node.Decorators {
		descriptorInfo, results := g.DescriptorGenerator.Generate(ctx, descriptor)
		if !results.IsEmpty() {
			return ReferencedTypeInfo{}, results
		}

		descriptors = append(descriptors, descriptorInfo)
	}

	typeInfo := ReferencedTypeInfo{
		Base:        baseType,
		Descriptors: descriptors,
	}

	return typeInfo, core.ResultList{}
}
