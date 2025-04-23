package aarch64translation_test

import (
	"testing"

	"alon.kr/x/aarch64codegen/registers"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"github.com/stretchr/testify/assert"
)

func TestValidNameToGPRegister(t *testing.T) {
	validNames := []struct {
		name     string
		register registers.GPRegister
	}{
		{"%x0", registers.GPRegisterX0},
		{"%x1", registers.GPRegisterX1},
		{"%x30", registers.GPRegisterX30},
		{"%xzr", registers.GPRegisterXZR},
	}

	for _, pair := range validNames {
		gpr, ok := aarch64translation.RegisterNameToAarch64GPRegister(pair.name)
		assert.True(t, ok)
		assert.Equal(t, pair.register, gpr)
	}
}

func TestInvalidNameToGPRegister(t *testing.T) {
	invalidNames := []string{"%x31", "%y0", "%", "%0", "%X0", "%x01", "x0", "x01"}
	for _, name := range invalidNames {
		_, ok := aarch64translation.RegisterNameToAarch64GPRegister(name)
		assert.False(t, ok)
	}
}
