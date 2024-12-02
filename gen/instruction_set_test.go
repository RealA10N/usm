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

type Instruction struct{}

type AddInstructionDefinition struct{}

func (AddInstructionDefinition) Names() []string {
	return []string{"ADD"}
}

func (AddInstructionDefinition) BuildInstruction(
	targets []*gen.RegisterInfo,
	arguments []*gen.ArgumentInfo,
) (Instruction, core.ResultList) {
	return Instruction{}, core.ResultList{}
}

func (AddInstructionDefinition) InferTargetTypes(
	ctx *gen.GenerationContext,
	targets []*gen.TypeInfo,
	arguments []*gen.TypeInfo,
) ([]*gen.TypeInfo, core.ResultList) {
	if len(arguments) != 2 {
		return nil, list.FromSingle(core.Result{{
			Type:    core.ErrorResult,
			Message: "expected exactly 2 arguments",
		}})
	}

	if len(targets) != 1 {
		return nil, list.FromSingle(core.Result{{
			Type:    core.ErrorResult,
			Message: "expected exactly 1 target",
		}})
	}

	if arguments[0] != arguments[1] {
		return nil, list.FromSingle(core.Result{{
			Type:    core.ErrorResult,
			Message: "expected both arguments to be of the same type",
		}})
	}

	return []*gen.TypeInfo{arguments[0]}, core.ResultList{}
}

func TestInstructionCreateTarget(t *testing.T) {
	isa := gen.NewInstructionSet([]gen.InstructionDefinition[Instruction]{
		&AddInstructionDefinition{},
	})

	src := core.NewSourceView("%c = ADD %a %b")
	tkns, err := lex.NewTokenizer().Tokenize(src)
	assert.NoError(t, err)

	// TODO: do no use parser here? test only the instruction set unit.
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

	_, results := isa.Build(ctx, node)
	assert.True(t, results.IsEmpty())

	target := registers.GetRegister("%c")
	assert.NotNil(t, target)
	assert.Equal(t, "%c", target.Name)
	assert.Equal(t, intType, target.Type)
	assert.Equal(t, src.Unmanaged().Subview(0, 2), target.Declaration)
}
