package aarch64managers

import (
	"math/big"

	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Aarch64TypeManager struct {
	// A lazy cache that contains all the integer types that have already been
	// declared and used.
	// This is a mapping of the full type name, e.g. "$64".
	IntegerTypes map[string]*gen.NamedTypeInfo
}

func (m *Aarch64TypeManager) createNewIntegerType(
	name string,
	size *big.Int,
) *gen.NamedTypeInfo {
	newType := gen.NewNamedTypeInfo(name, size, nil)
	m.IntegerTypes[name] = newType
	return newType
}

func (m *Aarch64TypeManager) GetType(name string) *gen.NamedTypeInfo {
	typ, ok := m.IntegerTypes[name]
	if ok {
		return typ
	}

	if size := aarch64translation.TypeNameToSize(name); size != nil {
		return m.createNewIntegerType(name, size)
	}

	return nil
}

func (m *Aarch64TypeManager) NewType(*gen.NamedTypeInfo) core.Result {
	return core.Result{
		{
			Type:    core.ErrorResult,
			Message: "Type declaration currently not supported in aarch64",
		},
	}
}

func NewTypeManager(*gen.GenerationContext) gen.TypeManager {
	return &Aarch64TypeManager{
		IntegerTypes: make(map[string]*gen.NamedTypeInfo),
	}
}
