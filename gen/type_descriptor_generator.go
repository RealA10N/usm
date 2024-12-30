package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type DescriptorGenerator struct{}

func NewDescriptorGenerator() FileContextGenerator[parse.TypeDecoratorNode, TypeDescriptorInfo] {
	return FileContextGenerator[parse.TypeDecoratorNode, TypeDescriptorInfo](
		&DescriptorGenerator{},
	)
}

// Valid type decorators should match the regex ".\d*" where the first rune is
// the decorator identifier (pointer, repeat, etc.), and immediately follows
// the an optional decimal number that is interpreted differently depending on
// decorator.
//
// This function parses the decorator string, and returns the decorator number,
// or an error if the decorator does not match the expected format. If a number
// is not provided, the default number is 1.
//
// Why don't we do this at the `parse` module? because the `parse` module parses
// the structure of tokens only, and does not look inside the content of the
// tokens. More specifically, it does not have access to the source context.
func (g *DescriptorGenerator) parseDescriptorAmount(
	ctx *FileGenerationContext,
	decorator parse.TypeDecoratorNode,
) (core.UsmUint, core.ResultList) {
	if decorator.Len() <= 1 {
		// 1 is the default amount for type decorators, when no explicit amount
		// is provided.
		return 1, core.ResultList{}
	}

	numView := decorator.Subview(1, decorator.Len())
	numStr := viewToSourceString(ctx, numView)
	num, err := core.ParseUint(numStr)

	if err != nil {
		return 0, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Failed to parse number in type decorator",
				Location: &numView,
			},
			{
				Type:    core.HintResult,
				Message: "Should be a positive, decimal number",
			},
		})
	}

	return num, core.ResultList{}
}

func (g *DescriptorGenerator) parsedDescriptorToGenDescriptorType(
	node parse.TypeDecoratorNode,
) (genType TypeDescriptorType, results core.ResultList) {
	switch node.Type {
	case parse.PointerTypeDecorator:
		return PointerTypeDescriptor, core.ResultList{}
	case parse.RepeatTypeDecorator:
		return RepeatTypeDescriptor, core.ResultList{}
	default:
		// notest
		return 0, list.FromSingle(core.Result{
			{
				Type:     core.InternalErrorResult,
				Message:  "Invalid type decorator",
				Location: &node.UnmanagedSourceView,
			},
		})
	}
}

func (g *DescriptorGenerator) Generate(
	ctx *FileGenerationContext,
	node parse.TypeDecoratorNode,
) (info TypeDescriptorInfo, results core.ResultList) {
	typ, results := g.parsedDescriptorToGenDescriptorType(node)
	if !results.IsEmpty() {
		return
	}

	amount, results := g.parseDescriptorAmount(ctx, node)
	if !results.IsEmpty() {
		return
	}

	return TypeDescriptorInfo{
		Type:   typ,
		Amount: amount,
	}, results
}
