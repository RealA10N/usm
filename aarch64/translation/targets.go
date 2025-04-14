package aarch64translation

import (
	"fmt"

	"alon.kr/x/aarch64codegen/registers"
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func AssertTargetsExactly(
	info *gen.InstructionInfo,
	count int,
) core.ResultList {
	if len(info.Targets) != count {
		return list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  fmt.Sprintf("Expected %d targets", count),
				Location: info.Declaration,
			},
		})
	}

	return core.ResultList{}
}

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
