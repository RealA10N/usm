// Converts AST nodes representing types (type declarations, etc.) into
// types metadata (calculating type structures, type sizes, etc.).

package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

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
func ParseDecoratorNum(
	genCtx *GenerationContext,
	dec parse.TypeDecoratorNode,
) (core.UsmUint, core.UsmError) {
	if dec.Len() <= 1 {
		return 1, nil
	}

	numView := dec.Subview(1, dec.Len())
	num, err := core.ParseUint(string(numView.Raw(genCtx.SourceContext)))

	if err != nil {
		return 0, core.GenericError{
			ErrorMessage:  "Failed to parse number in type decorator",
			HintMessage:   "Should be a positive, decimal number",
			ErrorLocation: numView,
		}
	}

	return num, nil
}

func CalculateTypeSizeFromTypeDecorators(
	genCtx *GenerationContext,
	baseTypeSize core.UsmUint,
	decorators []parse.TypeDecoratorNode,
) (core.UsmUint, core.UsmError) {
	if len(decorators) == 0 {
		return baseTypeSize, nil
	}

	topDecorator := decorators[len(decorators)-1]
	switch topDecorator.Type {

	case parse.PointerTypeDecorator:
		return genCtx.PointerSize, nil

	case parse.RepeatTypeDecorator:
		repeatNum, err := ParseDecoratorNum(genCtx, topDecorator)
		if err != nil {
			return 0, err
		}

		// TODO: optimization: iterative and not recursive
		return CalculateTypeSizeFromTypeDecorators(
			genCtx,
			repeatNum*baseTypeSize,
			decorators[:len(decorators)-1],
		)

	default:
		return 0, core.GenericError{
			ErrorMessage:  "Unknown type decorator",
			ErrorLocation: topDecorator.UnmanagedSourceView,
			IsInternal:    true,
		}
	}
}

func CalculateTypeSizeFromTypeNode(
	genCtx *GenerationContext,
	node parse.TypeNode,
) (core.UsmUint, core.UsmError) {
	typeName := string(node.Identifier.Raw(genCtx.SourceContext))
	typeInfo := genCtx.Types.GetType(typeName)

	if typeInfo == nil {
		return 0, core.GenericError{
			ErrorMessage:  "Undeclared type",
			ErrorLocation: node.Identifier,
		}
	}

	return CalculateTypeSizeFromTypeDecorators(
		genCtx,
		typeInfo.Size,
		node.Decorators,
	)
}

func CalculateTypeSizeFromTypeDeclaration(
	genCtx *GenerationContext,
	node parse.TypeDeclarationNode,
) (size core.UsmUint, err core.UsmError) {
	for _, node := range node.Fields.Nodes {
		var cur core.UsmUint
		cur, err = CalculateTypeSizeFromTypeNode(genCtx, node.Type)
		if err != nil {
			return
		}
		size += cur // TODO: possible int overflow?
	}

	return
}

// Converts an AST node representing a type declaration into a TypeInfo instance.
func TypeInfoFromTypeDeclaration(
	genCtx *GenerationContext,
	node parse.TypeDeclarationNode,
) (TypeInfo, core.UsmError) {
	size, err := CalculateTypeSizeFromTypeDeclaration(genCtx, node)
	return TypeInfo{
		Name: node.Identifier,
		Size: size,
	}, err
}
