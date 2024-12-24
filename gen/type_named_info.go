package gen

import "alon.kr/x/usm/core"

// A named type is a type that can has a distinct name.
// It either (1) a builtin type or (2) a type alias declared by the "type"
// keyword.
type NamedTypeInfo struct {
	Name string
	Size core.UsmUint

	// The source view of the type declaration.
	// Should be nil only if it is a builtin type.
	Declaration *core.UnmanagedSourceView
}
