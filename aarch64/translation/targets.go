package aarch64translation

import (
	"alon.kr/x/aarch64codegen/registers"
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func TargetToAarch64GPRegister(
	target gen.TargetInfo,
) (registers.GPRegister, core.ResultList) {
	regTarget, ok := target.(*gen.RegisterTargetInfo)
	if !ok {
		return 0, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected register target",
				Location: target.Declaration(),
			},
		})
	}

	return RegisterToAarch64GPRegister(regTarget.Register)
}

func TargetToAarch64GPorSPRegister(
	target gen.TargetInfo,
) (registers.GPorSPRegister, core.ResultList) {
	regTarget, ok := target.(*gen.RegisterTargetInfo)
	if !ok {
		return 0, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected register target",
				Location: target.Declaration(),
			},
		})
	}

	return RegisterToAarch64GPOrSPRegister(regTarget.Register)
}
