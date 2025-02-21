package aarch64translation

import (
	"strconv"

	"alon.kr/x/aarch64codegen/registers"
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func RegisterNameToAarch64GPRegister(
	name string,
) (registers.GPRegister, bool) {
	if len(name) < 2 {
		return 0, false
	}

	numStr := name[1:]
	num, err := strconv.ParseUint(numStr, 10, 64)
	gpr := registers.GPRegister(num)
	ok := name[0] == 'X' && err == nil && gpr.Validate() == nil
	return gpr, ok
}

func RegisterToAarch64GPRegister(
	register *gen.RegisterInfo,
) (registers.GPRegister, core.ResultList) {
	name := register.Name[1:]
	gpr, ok := RegisterNameToAarch64GPRegister(name)

	if !ok {
		return 0, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected Aarch64 General Purpose register",
				Location: &register.Declaration,
			},
		})
	}

	return gpr, core.ResultList{}
}
