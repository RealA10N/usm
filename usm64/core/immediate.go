package usm64core

import (
	"fmt"
	"math/big"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Immediate struct {
	value       uint64
	declaration core.UnmanagedSourceView
}

func NewImmediate(immediate gen.ImmediateInfo) (Immediate, core.ResultList) {
	if !immediate.Value.IsInt64() && !immediate.Value.IsUint64() {
		v := immediate.Declaration()
		return Immediate{}, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Immediate overflows 64 bits",
				Location: &v,
			},
		})
	}

	m := new(big.Int).Lsh(big.NewInt(1), 64)
	value := immediate.Value.Mod(immediate.Value, m).Uint64()

	return Immediate{
		value:       value,
		declaration: immediate.Declaration(),
	}, core.ResultList{}
}

func (i Immediate) Value(*EmulationContext) uint64 {
	return i.value
}

func (i Immediate) String(*EmulationContext) string {
	return fmt.Sprintf("#%d", i)
}

func (i Immediate) Declaration() core.UnmanagedSourceView {
	return i.declaration
}
