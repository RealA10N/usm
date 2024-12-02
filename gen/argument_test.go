package gen_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func TestUndefinedRegisterArgument(t *testing.T) {
	src := core.NewSourceView("%a")
	node := parse.RegisterNode{UnmanagedSourceView: src.Unmanaged()}

	ctx := gen.GenerationContext[Instruction]{
		ArchInfo:      gen.ArchInfo{PointerSize: 8},
		SourceContext: src.Ctx(),
		Types:         &TypeMap{},
		Registers:     &RegisterMap{},
		Instructions:  &InstructionMap{},
	}

	generator := gen.ArgumentGenerator[Instruction]{}
	_, results := generator.Generate(&ctx, node)

	assert.EqualValues(t, 1, results.Len())

	result := results.Head.Value
	assert.Equal(t, core.Result{{
		Type:     core.ErrorResult,
		Message:  "Undefined register used as argument",
		Location: &node.UnmanagedSourceView,
	}}, result)
}
