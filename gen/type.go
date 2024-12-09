// Converts AST nodes representing types (type declarations, etc.) into
// types metadata (calculating type structures, type sizes, etc.).

package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

// A named type is a type that can has a distinct name.
// It either (1) a builtin type or (2) a type alias declared by the "type"
// keyword.
type NamedTypeInfo struct {
	Name        string
	Size        core.UsmUint
	Declaration core.UnmanagedSourceView
}

type TypeDescriptorType uint8

const (
	PointerTypeDescriptor TypeDescriptorType = iota
	RepeatTypeDescriptor
)

type TypeDescriptorInfo struct {
	Type   TypeDescriptorType
	Amount core.UsmUint
}

// A referenced type is a combination of a basic type with (possibly zero)
// type decorators that wrap it.
// For example, if `$32“ is a basic named type, then `$32 *`, which is a
// pointer to that type is a referenced type with the `$32` named type as it's
// base type, and the pointer as a decorator.
type ReferencedTypeInfo struct {
	// A pointer to the base, named type that this type reference refers to.
	Base        *NamedTypeInfo
	Size        core.UsmUint
	Descriptors []TypeDescriptorInfo
}

func (info ReferencedTypeInfo) Equals(other ReferencedTypeInfo) bool {
	if info.Base != other.Base {
		return false
	}

	if info.Size != other.Size {
		return false
	}

	if len(info.Descriptors) != len(other.Descriptors) {
		return false
	}

	for i := range info.Descriptors {
		if info.Descriptors[i] != other.Descriptors[i] {
			return false
		}
	}

	return true
}

// MARK: Manager

type TypeManager interface {
	// Query a already seen before type, and get the type information if it
	// exists. Returns nil if the if a type with the provided name has not yet
	// been defined.
	//
	// The implementation should also return information about builtin types,
	// although the creation of such types can be possibly done lazily.
	GetType(name string) *NamedTypeInfo

	// Register a new type with the provided name and type information.
	// The generator will call this method when it encounters a new type
	// definition.
	//
	// The implementation should raise an error if the new registered type is
	// invalid. It can however assume that the type name is unique and has not
	// been defined before (GetType() returned nil on it).
	NewType(typ *NamedTypeInfo) core.Result
}

// MARK: Descriptor Generator

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

// MARK: Ref'ed Generator

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

func (g *ReferencedTypeGenerator) calculateTypeSize(
	ctx *FileGenerationContext,
	node parse.TypeNode,
	baseType *NamedTypeInfo,
	descriptors []TypeDescriptorInfo,
) (core.UsmUint, core.ResultList) {
	size := core.UsmUint(baseType.Size)

	for _, descriptor := range descriptors {
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

	size, results := g.calculateTypeSize(ctx, node, baseType, descriptors)
	if !results.IsEmpty() {
		return ReferencedTypeInfo{}, results
	}

	typeInfo := ReferencedTypeInfo{
		Base:        baseType,
		Size:        size,
		Descriptors: descriptors,
	}

	return typeInfo, core.ResultList{}
}

// MARK: Named Generator

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
				Location: &typeInfo.Declaration,
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

	referencedType := node.Fields.Nodes[0].Type
	referencedTypeInfo, results := g.ReferencedTypeGenerator.Generate(ctx, referencedType)
	if !results.IsEmpty() {
		return nil, results
	}

	typeInfo = &NamedTypeInfo{
		Name:        identifier,
		Size:        referencedTypeInfo.Size,
		Declaration: declaration,
	}

	result := ctx.Types.NewType(typeInfo)
	if result != nil {
		return nil, list.FromSingle(result)
	}

	return typeInfo, core.ResultList{}
}
