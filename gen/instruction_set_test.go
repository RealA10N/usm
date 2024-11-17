package gen_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

type AddInstructionDef struct{}

func (AddInstructionDef) HasSideEffects() bool { return false }

func TestInstructionSetNoErr(t *testing.T) {
	addInst := &AddInstructionDef{}

	instDef := gen.InstructionDef{
		Names: []string{"add", "ADD"},
		Builder: func(targets []parse.ParameterNode, arguments []parse.ArgumentNode) (gen.Instruction, error) {
			return addInst, nil
		},
	}

	set := gen.NewInstructionSet([]gen.InstructionDef{instDef})
	src := core.NewSourceView("add")

	inst, err := set.Build(src.Ctx(), parse.InstructionNode{
		Operator: src.Unmanaged(),
	})

	assert.NoError(t, err)
	assert.Same(t, addInst, inst)
}
