package gen

import "alon.kr/x/usm/core"

// MARK: Info Types

type TypeInfo struct {
	Name core.UnmanagedSourceView
	Size core.UsmUint
}

type RegisterInfo struct {
	// The name of the register, as it appears in the source code.
	Name core.UnmanagedSourceView

	// data oriented design: This is an index into the type table, of the type
	// of the register.
	TypeIndex core.UsmUint
}

type ArchInfo struct {
	PointerSize core.UsmUint // The size of a pointer in bytes.
}

// A context object that is initialized empty, but gets propagated and filled
// with information as the code generation process continues, while iterating
// over the AST nodes.
type GenerationContext struct {
	ArchInfo
	Types     []TypeInfo
	Registers []RegisterInfo

	// TODO: add globals, functions, etc.
}
