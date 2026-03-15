package aarch64isa_test

import (
	"fmt"
	"testing"

	"alon.kr/x/aarch64codegen/immediates"
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/aarch64codegen/registers"
	aarch64isa "alon.kr/x/usm/aarch64/isa"
)

func TestSubsExpectedCodegen(t *testing.T) {
	def := aarch64isa.NewSubs()

	testCases := []struct {
		src      string
		expected instructions.Instruction
	}{
		{
			"%xzr = subs %x0 %x1\n",
			instructions.NewSubsShiftedRegister(
				registers.GPRegisterXZR,
				registers.GPRegisterX0,
				registers.GPRegisterX1,
			),
		},
		{
			"%x0 = subs %x1 %x2\n",
			instructions.NewSubsShiftedRegister(
				registers.GPRegisterX0,
				registers.GPRegisterX1,
				registers.GPRegisterX2,
			),
		},
		{
			"%xzr = subs %x0 $12 #1234\n",
			instructions.NewSubsImmediate(
				registers.GPRegisterXZR,
				registers.GPorSPRegisterX0,
				immediates.Immediate12(1234),
			),
		},
		{
			"%x0 = subs %sp $12 #0\n",
			instructions.NewSubsImmediate(
				registers.GPRegisterX0,
				registers.GPorSPRegisterSP,
				immediates.Immediate12(0),
			),
		},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprint(idx), func(t *testing.T) {
			assertExpectedCodegen(t, def, testCase.expected, testCase.src)
		})
	}
}
