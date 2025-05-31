package aarch64translation

import (
	"alon.kr/x/aarch64codegen/registers"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func TargetToAarch64GPRegister(
	target *gen.TargetInfo,
) (registers.GPRegister, core.ResultList) {
	return RegisterToAarch64GPRegister(target.Register)
}

func TargetToAarch64GPorSPRegister(
	target *gen.TargetInfo,
) (registers.GPorSPRegister, core.ResultList) {
	return RegisterToAarch64GPOrSPRegister(target.Register)
}
