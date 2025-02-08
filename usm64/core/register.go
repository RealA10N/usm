package usm64core

import (
	"fmt"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Register struct {
	Name        string
	declaration *core.UnmanagedSourceView
}

func NewRegister(arg *gen.RegisterArgumentInfo) (Register, core.ResultList) {
	return Register{
		Name:        arg.Register.Name,
		declaration: arg.Declaration(),
	}, core.ResultList{}
}

func (r Register) Value(ctx *EmulationContext) uint64 {
	return ctx.Registers[r.Name]
}

func (r Register) String(ctx *EmulationContext) string {
	return fmt.Sprintf("%s (#%d)", r.Name, r.Value(ctx))
}

func (r Register) Declaration() *core.UnmanagedSourceView {
	return r.declaration
}
