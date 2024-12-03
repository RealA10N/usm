// Converts AST nodes representing types (type declarations, etc.) into
// types metadata (calculating type structures, type sizes, etc.).

package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type TypeInfo struct {
	Name        string
	Size        core.UsmUint
	Declaration core.UnmanagedSourceView
}

type TypeManager interface {
	// Query a already seen before type, and get the type information if it
	// exists. Returns nil if the if a type with the provided name has not yet
	// been defined.
	//
	// The implementation should also return information about builtin types,
	// although the creation of such types can be possibly done lazily.
	GetType(name string) *TypeInfo

	// Register a new type with the provided name and type information.
	// The generator will call this method when it encounters a new type
	// definition.
	//
	// The implementation should raise an error if the new registered type is
	// invalid, for example if there already exists a type with the same name,
	// or if its a builtin type.
	NewType(typ *TypeInfo) core.Result
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
func ParseDecoratorNum[InstT BaseInstruction](
	genCtx *GenerationContext[InstT],
	dec parse.TypeDecoratorNode,
) (core.UsmUint, core.Result) {
	if dec.Len() <= 1 {
		return 1, nil
	}

	numView := dec.Subview(1, dec.Len())
	num, err := core.ParseUint(string(numView.Raw(genCtx.SourceContext)))

	if err != nil {
		return 0, core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Failed to parse number in type decorator",
				Location: &numView,
			},
			{
				Type:    core.HintResult,
				Message: "Should be a positive, decimal number",
			},
		}
	}

	return num, nil
}

func CalculateTypeSizeFromTypeDecorators[InstT BaseInstruction](
	genCtx *GenerationContext[InstT],
	baseTypeSize core.UsmUint,
	decorators []parse.TypeDecoratorNode,
) (core.UsmUint, core.Result) {
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
		return 0, core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Unknown type decorator",
			Location: &topDecorator.UnmanagedSourceView,
		}}
	}
}

func CalculateTypeSizeFromTypeNode[InstT BaseInstruction](
	genCtx *GenerationContext[InstT],
	node parse.TypeNode,
) (core.UsmUint, core.Result) {
	typeName := string(node.Identifier.Raw(genCtx.SourceContext))
	typeInfo := genCtx.Types.GetType(typeName)

	if typeInfo == nil {
		return 0, NewUndefinedTypeResult(node.Identifier)
	}

	return CalculateTypeSizeFromTypeDecorators(
		genCtx,
		typeInfo.Size,
		node.Decorators,
	)
}

func CalculateTypeSizeFromTypeDeclaration[InstT BaseInstruction](
	genCtx *GenerationContext[InstT],
	node parse.TypeDeclarationNode,
) (size core.UsmUint, err core.Result) {
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
func TypeInfoFromTypeDeclaration[InstT BaseInstruction](
	genCtx *GenerationContext[InstT],
	node parse.TypeDeclarationNode,
) (TypeInfo, core.Result) {
	size, err := CalculateTypeSizeFromTypeDeclaration(genCtx, node)
	name := string(node.Identifier.Raw(genCtx.SourceContext))
	return TypeInfo{
		Name: name,
		Size: size,
	}, err
}
