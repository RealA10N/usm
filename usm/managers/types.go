package usmmanagers

import (
	"math/big"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// In the USM isa, there are infinitly many "base types", of the form "$<n>",
// each representing an integer of size n bits.
// We create those types lazily, and store them in the map.
type TypeMap map[string]*gen.NamedTypeInfo

// Check if the provided type name is of the form "$<n>" where n is a
// non-negative integer.
func (TypeMap) toBaseTypeSize(name string) *big.Int {
	// notice that there are multiple strings that can be mapped to the
	// same size, e.g. "$64" and "$064".

	size, ok := new(big.Int).SetString(name[1:], 10)
	if !ok || size.Sign() < 0 {
		return nil
	}

	return size
}

func (TypeMap) baseTypeSizeToCannonicalName(size *big.Int) string {
	return "$" + size.String()
}

func (m *TypeMap) GetType(name string) *gen.NamedTypeInfo {
	size := m.toBaseTypeSize(name)
	if size != nil {
		// is a base type, convert name to cannonical form
		name = m.baseTypeSizeToCannonicalName(size)
		if _, exists := (*m)[name]; !exists {
			// if does not exist, create it!
			(*m)[name] = gen.NewNamedTypeInfo(name, size, nil)
		}
	}

	return (*m)[name]
}

func (m *TypeMap) NewType(*gen.NamedTypeInfo) core.Result {
	return core.Result{
		{
			Type:    core.ErrorResult,
			Message: "Type declaration not supported yet in usm isa",
		},
	}
}

func NewTypeManager(*gen.GenerationContext) gen.TypeManager {
	return &TypeMap{}
}
