package aarch64managers

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// Aarch64 currently only supports the 64 bit type, named "$64".
type Aarch64TypeManager struct {
	BaseType *gen.NamedTypeInfo
}

func (m *Aarch64TypeManager) GetType(name string) *gen.NamedTypeInfo {
	if name == m.BaseType.Name {
		return m.BaseType
	} else {
		return nil
	}
}

func (m *Aarch64TypeManager) NewType(*gen.NamedTypeInfo) core.Result {
	return core.Result{
		{
			Type:    core.ErrorResult,
			Message: "Type declaration currently not supported in aarch64",
		},
	}
}

func NewTypeManager() gen.TypeManager {
	return &Aarch64TypeManager{
		BaseType: &gen.NamedTypeInfo{
			Name:        "$64",
			Size:        8,
			Declaration: nil,
		},
	}
}
