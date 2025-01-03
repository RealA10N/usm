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

// MARK: Add

type AddInstruction struct{}

func (i *AddInstruction) PossibleNextSteps() ([]gen.StepInfo, core.ResultList) {
	return []gen.StepInfo{gen.ContinueToNextInstruction{}}, core.ResultList{}
}

func (i *AddInstruction) String() string {
	return "ADD"
}

type AddInstructionDefinition struct{}

func (AddInstructionDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	return gen.BaseInstruction(&AddInstruction{}), core.ResultList{}
}

func (AddInstructionDefinition) InferTargetTypes(
	ctx *gen.FunctionGenerationContext,
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
	if !arguments[0].Equal(*arguments[1]) {
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

// MARK: Ret

type RetInstruction struct{}

func (i *RetInstruction) PossibleNextSteps() ([]gen.StepInfo, core.ResultList) {
	return []gen.StepInfo{gen.ReturnFromFunction{}}, core.ResultList{}
}

func (i *RetInstruction) String() string {
	return "RET"
}

type RetInstructionDefinition struct{}

func (RetInstructionDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	return gen.BaseInstruction(&RetInstruction{}), core.ResultList{}
}

func (RetInstructionDefinition) InferTargetTypes(
	ctx *gen.FunctionGenerationContext,
	targets []*gen.ReferencedTypeInfo,
	arguments []*gen.ReferencedTypeInfo,
) ([]gen.ReferencedTypeInfo, core.ResultList) {
	return []gen.ReferencedTypeInfo{}, core.ResultList{}
}

// MARK: Jump

type JumpInstruction struct {
	*gen.InstructionInfo
}

func (i *JumpInstruction) PossibleNextSteps() ([]gen.StepInfo, core.ResultList) {
	return []gen.StepInfo{gen.JumpToLabel{
		Label: i.Arguments[0].(*gen.LabelArgumentInfo).Label,
	}}, core.ResultList{}
}

func (i *JumpInstruction) OperatorString() string {
	return "JMP"
}

type JumpInstructionDefinition struct{}

func (JumpInstructionDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	return gen.BaseInstruction(&JumpInstruction{info}), core.ResultList{}
}

func (JumpInstructionDefinition) InferTargetTypes(
	ctx *gen.FunctionGenerationContext,
	targets []*gen.ReferencedTypeInfo,
	arguments []*gen.ReferencedTypeInfo,
) ([]gen.ReferencedTypeInfo, core.ResultList) {
	return []gen.ReferencedTypeInfo{}, core.ResultList{}
}

// MARK: Jump Zero

// JZ %condition .label
type JumpZeroInstruction struct {
	*gen.InstructionInfo
}

func (i *JumpZeroInstruction) PossibleNextSteps() ([]gen.StepInfo, core.ResultList) {
	return []gen.StepInfo{
		gen.JumpToLabel{Label: i.Arguments[1].(*gen.LabelArgumentInfo).Label},
		gen.ContinueToNextInstruction{},
	}, core.ResultList{}
}

func (i *JumpZeroInstruction) OperatorString() string {
	return "JZ"
}

type JumpZeroInstructionDefinition struct{}

func (JumpZeroInstructionDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	return gen.BaseInstruction(&JumpZeroInstruction{info}), core.ResultList{}
}

func (JumpZeroInstructionDefinition) InferTargetTypes(
	ctx *gen.FunctionGenerationContext,
	targets []*gen.ReferencedTypeInfo,
	arguments []*gen.ReferencedTypeInfo,
) ([]gen.ReferencedTypeInfo, core.ResultList) {
	return []gen.ReferencedTypeInfo{}, core.ResultList{}
}

// MARK: Instruction Map

type InstructionMap map[string]gen.InstructionDefinition

func (m *InstructionMap) GetInstructionDefinition(
	name string,
	node parse.InstructionNode,
) (gen.InstructionDefinition, core.ResultList) {
	instDef, ok := (*m)[name]
	if !ok {
		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "undefined instruction",
			Location: &node.Operator,
		}})
	}
	return instDef, core.ResultList{}
}

func TestInstructionCreateTarget(t *testing.T) {
	src := core.NewSourceView("%c = ADD %a %b\n")
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

	ctx := &gen.FunctionGenerationContext{
		FileGenerationContext: &gen.FileGenerationContext{
			GenerationContext: &testGenerationContext,
			SourceContext:     src.Ctx(),
			Types:             &types,
		},
		Registers: &registers,
	}

	generator := gen.NewInstructionGenerator()
	_, results := generator.Generate(ctx, node)
	assert.True(t, results.IsEmpty())

	target := registers.GetRegister("%c")
	assert.NotNil(t, target)
	assert.Equal(t, "%c", target.Name)
	assert.Equal(t, intType, target.Type.Base)
	assert.Equal(t, src.Unmanaged().Subview(0, 2), target.Declaration)
}
