package gen_test

import (
	"testing"

	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"
	"github.com/stretchr/testify/assert"
)

func TestInstructionSetNoErr(t *testing.T) {
	type AddInstructionDef struct{}
	addInst := &AddInstructionDef{}

	instDef := gen.InstructionDef{
		Names: []string{"add", "ADD"},
		Builder: func(targets []parse.ParameterNode, arguments []parse.ArgumentNode) (gen.Instruction, error) {
			return addInst, nil
		},
	}

	set := gen.NewInstructionSet([]gen.InstructionDef{instDef})
	src := source.NewSourceView("add")

	inst, err := set.Build(src.Ctx(), parse.InstructionNode{
		Operator: src.Unmanaged(),
	})

	assert.NoError(t, err)
	assert.Same(t, addInst, inst)
}
