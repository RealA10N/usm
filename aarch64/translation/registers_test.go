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
		{"X0", registers.X0},
		{"X1", registers.X1},
		{"X30", registers.X30},
	}

	for _, pair := range validNames {
		gpr, ok := aarch64translation.RegisterNameToAarch64GPRegister(pair.name)
		assert.True(t, ok)
		assert.Equal(t, pair.register, gpr)
	}
}

func TestInvalidNameToGPRegister(t *testing.T) {
	invalidNames := []string{"X31", "Y0", "", "0", "x0"}
	for _, name := range invalidNames {
		_, ok := aarch64translation.RegisterNameToAarch64GPRegister(name)
		assert.False(t, ok)
	}
}
