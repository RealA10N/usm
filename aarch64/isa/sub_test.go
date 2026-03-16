package aarch64isa_test

import (
	"fmt"
	"testing"

	"alon.kr/x/aarch64codegen/immediates"
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/aarch64codegen/registers"
	aarch64isa "alon.kr/x/usm/aarch64/isa"
)

func TestSubExpectedCodegen(t *testing.T) {
	def := aarch64isa.NewSub()

	testCases := []struct {
		src      string
		expected instructions.Instruction
	}{
		{
			"%x0 = sub %x1 %x2\n",
			instructions.NewSubShiftedRegister(
				registers.GPRegisterX0,
				registers.GPRegisterX1,
				registers.GPRegisterX2,
			),
		},
		{
			"%xzr = sub %xzr %xzr\n",
			instructions.NewSubShiftedRegister(
				registers.GPRegisterXZR,
				registers.GPRegisterXZR,
				registers.GPRegisterXZR,
			),
		},
		{
			"%x0 = sub %x1 $12 #1234\n",
			instructions.NewSubImmediate(
				registers.GPRegisterX0,
				registers.GPorSPRegisterX1,
				immediates.Immediate12(1234),
			),
		},
		{
			"%xzr = sub %sp $12 #0\n",
			instructions.NewSubImmediate(
				registers.GPRegisterXZR,
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
