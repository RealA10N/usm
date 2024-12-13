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

func (AddInstructionDefinition) BuildInstruction(
	targets []*gen.RegisterArgumentInfo,
	arguments []gen.ArgumentInfo,
) (Instruction, core.ResultList) {
	return Instruction{}, core.ResultList{}
}

func (AddInstructionDefinition) InferTargetTypes(
	ctx *gen.FunctionGenerationContext[Instruction],
	targets []*gen.ReferencedTypeInfo,
	arguments []*gen.ReferencedTypeInfo,
) ([]gen.ReferencedTypeInfo, core.ResultList) {
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

	// TODO: possible panic?
	if !arguments[0].Equals(*arguments[1]) {
		return nil, list.FromSingle(core.Result{{
			Type:    core.ErrorResult,
			Message: "expected both arguments to be of the same type",
		}})
	}

	return []gen.ReferencedTypeInfo{
		{
			Base:        arguments[0].Base,
			Descriptors: arguments[0].Descriptors,
		},
	}, core.ResultList{}
}

type InstructionMap map[string]gen.InstructionDefinition[Instruction]

func (m *InstructionMap) GetInstructionDefinition(
	name string,
) (gen.InstructionDefinition[Instruction], core.ResultList) {
	instDef, ok := (*m)[name]
	if !ok {
		return nil, list.FromSingle(core.Result{{
			Type:    core.ErrorResult,
			Message: "undefined instruction",
		}})
	}
	return instDef, core.ResultList{}
}

func TestInstructionCreateTarget(t *testing.T) {
	src := core.NewSourceView("%c = ADD %a %b")
	tkns, err := lex.NewTokenizer().Tokenize(src)
	assert.NoError(t, err)

	// TODO: do no use parser here? test only the instruction set unit.
	tknView := parse.NewTokenView(tkns)
	node, result := parse.NewInstructionParser().Parse(&tknView)
	assert.Nil(t, result)

	intType := &gen.NamedTypeInfo{Name: "$32", Size: 4}
	types := TypeMap{intType.Name: intType}

	intTypeRef := gen.ReferencedTypeInfo{
		Base: intType,
	}

	registers := RegisterMap{
		"%a": &gen.RegisterInfo{Name: "%a", Type: intTypeRef},
		"%b": &gen.RegisterInfo{Name: "%b", Type: intTypeRef},
	}

	ctx := &gen.FunctionGenerationContext[Instruction]{
		FileGenerationContext: &gen.FileGenerationContext[Instruction]{
			GenerationContext: &testGenerationContext,
			SourceContext:     src.Ctx(),
			Types:             &types,
		},
		Registers: &registers,
	}

	generator := gen.NewInstructionGenerator[Instruction]()
	_, results := generator.Generate(ctx, node)
	assert.True(t, results.IsEmpty())

	target := registers.GetRegister("%c")
	assert.NotNil(t, target)
	assert.Equal(t, "%c", target.Name)
	assert.Equal(t, intType, target.Type.Base)
	assert.Equal(t, src.Unmanaged().Subview(0, 2), target.Declaration)
}
