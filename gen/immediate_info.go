package gen

import (
	"math/big"

	"alon.kr/x/usm/core"
)

type ImmediateInfo struct {
	Type        ReferencedTypeInfo
	Value       *big.Int // TODO: Add floating types
	declaration core.UnmanagedSourceView
	// TODO: more complex and complete representation of immediate structs.
}

func (i *ImmediateInfo) GetType() *ReferencedTypeInfo {
	return &i.Type
}

func (i *ImmediateInfo) Declaration() core.UnmanagedSourceView {
	return i.declaration
}

func (i *ImmediateInfo) String() string {
	return i.Type.String() + " #" + i.Value.String()
}
