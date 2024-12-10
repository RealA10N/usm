package core_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"github.com/stretchr/testify/assert"
)

func TestParseUintDecimal(t *testing.T) {
	n, err := core.ParseUint("123")
	assert.Nil(t, err)
	assert.Equal(t, core.UsmUint(123), n)
}

// Currently, the expected behavior is to not support other bases.
// Parsing currently fails if provided no decimal digit runes (this is the
// expected behavior).
// TODO: consider adding support for other bases?

func TestParseUintBinary(t *testing.T) {
	n, err := core.ParseUint("0b101")
	assert.Nil(t, err)
	assert.Equal(t, core.UsmUint(5), n)
}

func TestParseUintPadding(t *testing.T) {
	_, err := core.ParseUint(" 123   ")
	assert.NotNil(t, err)
}
