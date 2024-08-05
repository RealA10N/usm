package source_test

import (
	"strings"
	"testing"

	"alon.kr/x/usm/source"

	"github.com/stretchr/testify/assert"
)

func TestReadSourceSimpleCase(t *testing.T) {
	data := "hello, world!"
	reader := strings.NewReader(data)
	view, err := source.ReadSource(reader)
	assert.NoError(t, err)
	assert.Equal(t, data, string(view.Raw()))
}
