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
	unmanaged := src.Unmanaged()
	node := parse.RegisterNode{UnmanagedSourceView: unmanaged}

	ctx := gen.GenerationContext[Instruction]{
		ArchInfo:      gen.ArchInfo{PointerSize: 8},
		SourceContext: src.Ctx(),
		Types:         &TypeMap{},
		Registers:     &RegisterMap{},
		Instructions:  &InstructionMap{},
	}

	generator := gen.NewArgumentGenerator[Instruction]()
	_, results := generator.Generate(&ctx, node)

	assert.EqualValues(t, 1, results.Len())

	result := results.Head.Value
	assert.Equal(t, gen.NewUndefinedRegisterResult(unmanaged), result)
}
