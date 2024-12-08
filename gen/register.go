package gen

import "alon.kr/x/usm/core"

// MARK: Info

type RegisterInfo struct {
	// The name of the register, as it appears in the source code.
	Name string

	// A pointer to the TypeInfo instance that corresponds to the type of the
	// register.
	Type *NamedTypeInfo

	// The first location in the source code in which the register is declared
	// or assigned a value.
	Declaration core.UnmanagedSourceView
}

// MARK: Manager

type RegisterManager interface {
	GetRegister(name string) *RegisterInfo
	NewRegister(reg *RegisterInfo) core.Result
}
