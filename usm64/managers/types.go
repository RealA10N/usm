package managers

import (
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

func NewTypeManager() gen.TypeManager {
	return &TypeMap{
		BaseType: &gen.NamedTypeInfo{
			Name:        "$64",
			Size:        8,
			Declaration: nil,
		},
	}
}
