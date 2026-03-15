package aarch64isa_test

import (
	"fmt"
	"testing"

	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/aarch64codegen/registers"
	aarch64isa "alon.kr/x/usm/aarch64/isa"
)

func TestRetExpectedCodegen(t *testing.T) {
	def := aarch64isa.NewRet()

	testCases := []struct {
		src      string
		expected instructions.Instruction
	}{
		{
			"ret\n",
			instructions.RET(registers.GPRegisterX30),
		},
		{
			"ret %x30\n",
			instructions.RET(registers.GPRegisterX30),
		},
		{
			"ret %x0\n",
			instructions.RET(registers.GPRegisterX0),
		},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprint(idx), func(t *testing.T) {
			assertExpectedCodegen(t, def, testCase.expected, testCase.src)
		})
	}
}
