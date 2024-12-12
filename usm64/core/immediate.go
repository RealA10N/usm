package usm64core

import (
	"math/big"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Immediate uint64

func NewImmediate(immediate gen.ImmediateInfo) (Immediate, core.ResultList) {
	if !immediate.Value.IsInt64() && !immediate.Value.IsUint64() {
		return Immediate(0), list.FromSingle(core.Result{
			{
				Type:    core.ErrorResult,
				Message: "Immediate overflows 64 bits",
			},
		})
	}

	m := new(big.Int).Lsh(big.NewInt(1), 64)
	value := immediate.Value.Mod(immediate.Value, m).Uint64()
	return Immediate(value), core.ResultList{}
}

func (i Immediate) Value(*EmulationContext) uint64 {
	return uint64(i)
}
