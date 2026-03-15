package aarch64isa_test

import (
	"fmt"
	"testing"

	"alon.kr/x/aarch64codegen/immediates"
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/aarch64codegen/registers"
	aarch64isa "alon.kr/x/usm/aarch64/isa"
)

func TestMovzExpectedCodegen(t *testing.T) {
	def := aarch64isa.NewMovz()

	testCases := []struct {
		src      string
		expected instructions.Instruction
	}{
		{
			"%x0 = movz $16 #0xffff\n",
			instructions.MOVZ(
				registers.GPRegisterX0,
				immediates.Immediate16(0xffff),
				instructions.MovShift0,
			),
		},
		{
			"%x1 = movz $16 #0\n",
			instructions.MOVZ(
				registers.GPRegisterX1,
				immediates.Immediate16(0),
				instructions.MovShift0,
			),
		},
		{
			"%x2 = movz $16 #0x1234 $8 #16\n",
			instructions.MOVZ(
				registers.GPRegisterX2,
				immediates.Immediate16(0x1234),
				instructions.MovShift16,
			),
		},
		{
			"%xzr = movz $16 #0xabcd $8 #32\n",
			instructions.MOVZ(
				registers.GPRegisterXZR,
				immediates.Immediate16(0xabcd),
				instructions.MovShift32,
			),
		},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprint(idx), func(t *testing.T) {
			assertExpectedCodegen(t, def, testCase.expected, testCase.src)
		})
	}
}
