package gen

import "alon.kr/x/usm/core"

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
