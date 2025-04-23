package managers

import (
	"math/big"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type TypeMap struct {
	BaseType *gen.NamedTypeInfo
}

func (m *TypeMap) GetType(name string) *gen.NamedTypeInfo {
	if name == m.BaseType.Name {
		return m.BaseType
	} else {
		return nil
	}
}

func (m *TypeMap) NewType(*gen.NamedTypeInfo) core.Result {
	return core.Result{
		{
			Type:    core.ErrorResult,
			Message: "Type declaration not supported in usm64",
		},
	}
}

func NewTypeManager(*gen.GenerationContext) gen.TypeManager {
	return &TypeMap{
		BaseType: gen.NewNamedTypeInfo("$64", big.NewInt(64), nil),
	}
}
