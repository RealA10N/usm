package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalMap(t *testing.T) {
	gm := NewGlobalMap(nil)
	assert.Empty(t, gm.GetAllGlobals())

	info := &FunctionInfo{Name: "foo"}
	global := NewFunctionGlobalInfo(info)

	results := gm.NewGlobal(global)
	assert.True(t, results.IsEmpty())

	gotGlobal := gm.GetGlobal("foo")
	assert.Equal(t, global, gotGlobal)
	assert.Equal(t, info.Name, gotGlobal.Name())
}
