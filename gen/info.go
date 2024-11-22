package gen

import (
	"alon.kr/x/usm/core"
)

// MARK: Info Types

type TypeInfo struct {
	Name core.UnmanagedSourceView
	Size core.UsmUint
}

type RegisterInfo struct {
	// The name of the register, as it appears in the source code.
	Name core.UnmanagedSourceView

	// A pointer to the TypeInfo instance that corresponds to the type of the
	// register.
	Type *TypeInfo
}

type ArchInfo struct {
	PointerSize core.UsmUint // The size of a pointer in bytes.
}

// MARK: Managers

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
	RegisterType(name string, typ *TypeInfo) core.Result
}

// MARK: Generation Context

// A context object that is initialized empty, but gets propagated and filled
// with information as the code generation process continues, while iterating
// over the AST nodes.
type GenerationContext struct {
	ArchInfo
	core.SourceContext

	Types TypeManager
	// TODO: add registers info.

	// TODO: add globals, functions, etc.
}
