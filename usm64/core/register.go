package usm64core

import (
	"fmt"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// MARK: Register

type Register string

func NewRegister(register *gen.RegisterInfo) (Register, core.ResultList) {
	return Register(register.Name), core.ResultList{}
}

func (r Register) Value(ctx *EmulationContext) uint64 {
	return ctx.Registers[r]
}

func (r Register) String(ctx *EmulationContext) string {
	return fmt.Sprintf("%s (#%d)", r, r.Value(ctx))
}
