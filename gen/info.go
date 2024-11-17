package gen

import "alon.kr/x/usm/core"

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

// A context object that is initialized empty, but gets propagated and filled
// with information as the code generation process continues, while iterating
// over the AST nodes.
type FileInfo struct {
	Types     []TypeInfo
	Registers []RegisterInfo

	// TODO: add globals, functions, etc.
}
