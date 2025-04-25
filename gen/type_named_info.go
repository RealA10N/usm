package gen

import (
	"math/big"

	"alon.kr/x/usm/core"
)

// A named type is a type that can has a distinct name.
// It either (1) a builtin type or (2) a type alias declared by the "type"
// keyword.
type NamedTypeInfo struct {
	// The full name of the type, including the "$" prefix.
	Name string

	// The size of the type in bits.
	Size *big.Int

	// The source view of the type declaration.
	// Should be nil only if it is a builtin type.
	Declaration *core.UnmanagedSourceView
}

func NewNamedTypeInfo(
	name string,
	size *big.Int,
	declaration *core.UnmanagedSourceView,
) *NamedTypeInfo {
	return &NamedTypeInfo{
		Name:        name,
		Size:        size,
		Declaration: declaration,
	}
}

func (n NamedTypeInfo) String() string {
	return n.Name
}
