package gen_test

import (
	"testing"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

type AddInstructionDefinition struct{}

func (AddInstructionDefinition) Names() []string {
	return []string{"ADD"}
}

type AddInstruction struct {
	Left   string
	Right  string
	Target string
}

func (AddInstructionDefinition) BuildInstruction(
	targets []*gen.RegisterInfo,
	arguments []*gen.ArgumentInfo,
) (gen.Instruction, core.ResultList) {
	return AddInstruction{
		Left:   arguments[0].Type.Name,
		Right:  arguments[1].Type.Name,
		Target: targets[0].Name,
	}, core.ResultList{}
}

func (AddInstructionDefinition) InferTargetTypes(
	ctx *gen.GenerationContext,
	targets []*gen.TypeInfo,
	arguments []*gen.TypeInfo,
) ([]*gen.TypeInfo, core.ResultList) {
	if len(arguments) != 2 {
		return nil, list.FromSlice([]core.Result{
			&core.GenericResult{
				Type:    core.ErrorResult,
				Message: "expected exactly 2 arguments",
			},
		})
	}

	if len(targets) != 1 {
		return nil, list.FromSlice([]core.Result{
			&core.GenericResult{
				Type:    core.ErrorResult,
				Message: "expected exactly 1 target",
			},
		})
	}

	if arguments[0] != arguments[1] {
		return nil, list.FromSlice([]core.Result{
			&core.GenericResult{
				Type:    core.ErrorResult,
				Message: "expected both arguments to be of the same type",
			},
		})
	}

	return []*gen.TypeInfo{arguments[0]}, core.ResultList{}
}

func TestInstructionSetNoErr(t *testing.T) {
	isa := gen.NewInstructionSet([]gen.InstructionDefinition{
		&AddInstructionDefinition{},
	})

	src := core.NewSourceView("%c = ADD %a %b")
	tkns, err := lex.NewTokenizer().Tokenize(src)
	assert.NoError(t, err)

	tknView := parse.NewTokenView(tkns)
	node, result := parse.NewInstructionParser().Parse(&tknView)
	assert.Nil(t, result)

	intType := &gen.TypeInfo{Name: "$32", Size: 4}
	types := TypeMap{intType.Name: intType}

	registers := RegisterMap{
		"%a": &gen.RegisterInfo{Name: "%a", Type: intType},
		"%b": &gen.RegisterInfo{Name: "%b", Type: intType},
	}

	ctx := &gen.GenerationContext{
		ArchInfo:      gen.ArchInfo{PointerSize: 8},
		SourceContext: src.Ctx(),
		Types:         &types,
		Registers:     &registers,
	}

	inst, results := isa.Build(ctx, node)
	assert.True(t, results.IsEmpty())

	addInst := inst.(AddInstruction)
	assert.Equal(t, "$32", addInst.Left)
	assert.Equal(t, "$32", addInst.Right)
	assert.Equal(t, "%c", addInst.Target)
}
