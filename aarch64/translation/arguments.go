package aarch64translation

import (
	"math/big"

	"alon.kr/x/aarch64codegen/immediates"
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/aarch64codegen/registers"
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func ArgumentToAarch64GPRegister(
	argument gen.ArgumentInfo,
) (registers.GPRegister, core.ResultList) {
	register, ok := argument.(*gen.RegisterArgumentInfo)
	if !ok {
		return 0, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected register argument",
				Location: argument.Declaration(),
			},
		})
	}

	return RegisterToAarch64GPRegister(register.Register)
}

func ArgumentToAarch64GPorSPRegister(
	argument gen.ArgumentInfo,
) (registers.GPorSPRegister, core.ResultList) {
	register, ok := argument.(*gen.RegisterArgumentInfo)
	if !ok {
		return 0, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected register argument",
				Location: argument.Declaration(),
			},
		})
	}
	return RegisterToAarch64GPOrSPRegister(register.Register)
}

func ArgumentToImmediateInfo(
	argument gen.ArgumentInfo,
) (*gen.ImmediateInfo, core.ResultList) {
	imm, ok := argument.(*gen.ImmediateInfo)
	if !ok {
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected immediate argument",
				Location: argument.Declaration(),
			},
		})
	}

	return imm, core.ResultList{}
}

func ArgumentToAarch64Immediate12(
	argument gen.ArgumentInfo,
) (immediates.Immediate12, core.ResultList) {
	// TODO: remove code duplication with other immediate types.

	info, results := ArgumentToImmediateInfo(argument)
	if !results.IsEmpty() {
		return 0, results
	}

	results = AssertIntegerTypeOfSize(info.Type, big.NewInt(12))
	if !results.IsEmpty() {
		return 0, results
	}

	return BigIntToAarch64Immediate12(argument.Declaration(), info.Value)
}

func BigIntToAarch64Immediate12(
	view *core.UnmanagedSourceView,
	bigInt *big.Int,
) (value immediates.Immediate12, results core.ResultList) {
	isInvalid := bigInt.Sign() < 0 || bigInt.BitLen() > 12
	if isInvalid {
		goto fail
	}

	value = immediates.Immediate12(bigInt.Uint64())
	if value.Validate() != nil {
		goto fail
	}

	return

fail:
	return 0, list.FromSingle(core.Result{
		{
			Type:     core.ErrorResult,
			Message:  "Expected 12-bit unsigned integer",
			Location: view,
		},
	})
}

func ArgumentToAarch64Immediate16(
	argument gen.ArgumentInfo,
) (immediates.Immediate16, core.ResultList) {
	info, results := ArgumentToImmediateInfo(argument)
	if !results.IsEmpty() {
		return 0, results
	}

	results = AssertIntegerTypeOfSize(info.Type, big.NewInt(16))
	if !results.IsEmpty() {
		return 0, results
	}

	return BigIntToAarch64Immediate16(argument.Declaration(), info.Value)
}

func BigIntToAarch64Immediate16(
	view *core.UnmanagedSourceView,
	bigInt *big.Int,
) (value immediates.Immediate16, results core.ResultList) {
	isInvalid := bigInt.Sign() < 0 || bigInt.BitLen() > 16
	if isInvalid {
		goto fail
	}

	value = immediates.Immediate16(bigInt.Uint64())
	if value.Validate() != nil {
		goto fail
	}

	return

fail:
	return 0, list.FromSingle(core.Result{
		{
			Type:     core.ErrorResult,
			Message:  "Expected 16-bit unsigned integer",
			Location: view,
		},
	})
}

func BigIntToAarch64MovShift(
	view *core.UnmanagedSourceView,
	bigInt *big.Int,
) (instructions.MovShift, core.ResultList) {
	value, results := gen.AssertBigIntInSet(view, bigInt, []int64{0, 16, 32, 48})
	if !results.IsEmpty() {
		return 0, results
	}

	shift := instructions.MovShift(value / 16)
	if shift.Validate() != nil {
		return 0, list.FromSingle(core.Result{
			{
				Type:     core.InternalErrorResult,
				Message:  "Unexpected shift value",
				Location: view,
			},
		})
	}

	return shift, core.ResultList{}
}

func ArgumentToAarch64MovShift(
	argument gen.ArgumentInfo,
) (instructions.MovShift, core.ResultList) {
	info, results := ArgumentToImmediateInfo(argument)
	if !results.IsEmpty() {
		return 0, results
	}

	return BigIntToAarch64MovShift(argument.Declaration(), info.Value)
}

func ArgumentToLabelInfo(argument gen.ArgumentInfo) (*gen.LabelInfo, core.ResultList) {
	label, ok := argument.(*gen.LabelArgumentInfo)
	if !ok {
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected label argument",
				Location: argument.Declaration(),
			},
		})
	}

	return label.Label, core.ResultList{}
}

func ArgumentToFunctionInfo(
	argument gen.ArgumentInfo,
) (*gen.FunctionInfo, core.ResultList) {
	globalArg, ok := argument.(*gen.GlobalArgumentInfo)
	if !ok {
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected global argument",
				Location: argument.Declaration(),
			},
		})
	}

	info, ok := globalArg.GlobalInfo.(*gen.FunctionGlobalInfo)
	if !ok {
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected function global argument",
				Location: argument.Declaration(),
			},
		})
	}

	return info.FunctionInfo, core.ResultList{}
}
