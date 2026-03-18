package aarch64translation

import (
	"alon.kr/x/aarch64codegen/registers"
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func TargetToAarch64GPRegister(
	target gen.ArgumentInfo,
) (registers.GPRegister, core.ResultList) {
	regArg, ok := target.(*gen.RegisterArgumentInfo)
	if !ok {
		v := target.Declaration()
		return 0, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Expected a register target",
			Location: v,
		}})
	}

	return RegisterToAarch64GPRegister(regArg.Register)
}

func TargetToAarch64GPorSPRegister(
	target gen.ArgumentInfo,
) (registers.GPorSPRegister, core.ResultList) {
	regArg, ok := target.(*gen.RegisterArgumentInfo)
	if !ok {
		v := target.Declaration()
		return 0, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Expected a register target",
			Location: v,
		}})
	}

	return RegisterToAarch64GPOrSPRegister(regArg.Register)
}
