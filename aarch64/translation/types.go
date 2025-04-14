package aarch64translation

import (
	"fmt"
	"math/big"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// TypeNameToSize converts a type name to its size in bits.
// For example, "$64" will return 64, and "$128" will return 128.
// If the type name is not a valid integer type name, it returns nil.
func TypeNameToSize(name string) *big.Int {
	withoutPrefix := name[1:]
	value, ok := new(big.Int).SetString(withoutPrefix, 10)
	if !ok {
		return nil
	}
	return value
}

// Checks if the provided type is an integer type.
// If it is, returns the size of the integer type, in bits.
// Otherwise, returns nil.
func IsIntegerType(typ gen.ReferencedTypeInfo) *big.Int {
	if !typ.IsPure() {
		return nil
	}

	return TypeNameToSize(typ.Base.Name)
}

func AssertIntegerTypeOfSize(
	typ gen.ReferencedTypeInfo,
	expectedSize *big.Int,
) core.ResultList {
	actualSize := IsIntegerType(typ)

	if actualSize == nil || actualSize.Cmp(expectedSize) != 0 {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  fmt.Sprintf("Expected integer type $%s", expectedSize.String()),
				Location: typ.Declaration,
			},
		})
	}

	return core.ResultList{}
}
