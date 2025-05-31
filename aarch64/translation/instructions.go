package aarch64translation

import (
	"alon.kr/x/aarch64codegen/immediates"
	"alon.kr/x/aarch64codegen/registers"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func ValidateBinaryInstruction(
	info *gen.InstructionInfo,
) core.ResultList {
	results := gen.AssertTargetsExactly(info, 1)

	argumentResults := gen.AssertArgumentsExactly(info, 2)
	results.Extend(&argumentResults)

	return results
}

// BinaryInstructionToAarch64 converts a binary instruction to its AArch64
// representation.
//
// Assumes that the number of targets and arguments has already been validated
// using the ValidateBinaryInstruction function.
func BinaryInstructionToAarch64(
	info *gen.InstructionInfo,
) (Xd, Xn, Xm registers.GPRegister, results core.ResultList) {
	Xd, curResults := TargetToAarch64GPRegister(info.Targets[0])
	results.Extend(&curResults)

	Xn, curResults = ArgumentToAarch64GPRegister(info.Arguments[0])
	results.Extend(&curResults)

	Xm, curResults = ArgumentToAarch64GPRegister(info.Arguments[1])
	results.Extend(&curResults)

	return
}

func Immediate12InstructionToAarch64(
	info *gen.InstructionInfo,
) (
	Xd, Xn registers.GPorSPRegister,
	imm immediates.Immediate12,
	results core.ResultList,
) {
	Xd, curResults := TargetToAarch64GPorSPRegister(info.Targets[0])
	results.Extend(&results)

	Xn, curResults = ArgumentToAarch64GPorSPRegister(info.Arguments[0])
	results.Extend(&curResults)

	// TODO: Add shifted immediate support
	imm, curResults = ArgumentToAarch64Immediate12(info.Arguments[1])
	results.Extend(&curResults)

	return
}

// Immediate12GPRegisterTargetInstructionToAarch64 converts a binary instruction
// with its target as GPRegister, first argument as a GPorSPRegister and its
// second argument as an Immediate12 to it's codegen representation.
//
// Assumes that the number of targets and arguments has already been validated
// using the ValidateBinaryInstruction function.
func Immediate12GPRegisterTargetInstructionToAarch64(
	info *gen.InstructionInfo,
) (
	Xd registers.GPRegister,
	Xn registers.GPorSPRegister,
	imm immediates.Immediate12,
	results core.ResultList,
) {
	Xd, curResults := TargetToAarch64GPRegister(info.Targets[0])
	results.Extend(&curResults)

	Xn, curResults = ArgumentToAarch64GPorSPRegister(info.Arguments[0])
	results.Extend(&curResults)

	imm, curResults = ArgumentToAarch64Immediate12(info.Arguments[1])
	results.Extend(&curResults)

	return
}
