package gen_test

import (
	"math/big"
	"testing"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

// MARK: Add

type Add struct {
	gen.NonBranchingInstruction
}

func NewAdd() gen.InstructionDefinition {
	return Add{}
}

func (Add) Operator(*gen.InstructionInfo) string {
	return "add"
}

func (Add) Validate(info *gen.InstructionInfo) core.ResultList {
	return core.ResultList{}
}

// MARK: Ret

type Ret struct{}

func NewRet() gen.InstructionDefinition {
	return Ret{}
}

func (Ret) Operator(*gen.InstructionInfo) string {
	return "ret"
}
func (Ret) PossibleNextSteps(*gen.InstructionInfo) (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleReturn: true}, core.ResultList{}
}

func (Ret) Validate(*gen.InstructionInfo) core.ResultList {
	return core.ResultList{}
}

// MARK: Jump

type Jump struct{}

func NewJump() gen.InstructionDefinition {
	return Jump{}
}

func (Jump) PossibleNextSteps(i *gen.InstructionInfo) (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{
		PossibleBranches: []*gen.LabelInfo{
			i.Arguments[0].(*gen.LabelArgumentInfo).Label,
		},
	}, core.ResultList{}
}

func (Jump) Operator(*gen.InstructionInfo) string {
	return "j"
}

func (Jump) Validate(info *gen.InstructionInfo) core.ResultList {
	return core.ResultList{}
}

// MARK: Jump Zero

// JZ %condition .label
type JumpZero struct{}

func NewJumpZero() gen.InstructionDefinition {
	return JumpZero{}
}

func (JumpZero) Operator(*gen.InstructionInfo) string {
	return "jz"
}

func (JumpZero) PossibleNextSteps(i *gen.InstructionInfo) (gen.StepInfo, core.ResultList) {
	label := i.Arguments[1].(*gen.LabelArgumentInfo).Label
	return gen.StepInfo{
		PossibleBranches: []*gen.LabelInfo{label},
		PossibleContinue: true,
	}, core.ResultList{}
}

func (JumpZero) Validate(info *gen.InstructionInfo) core.ResultList {
	return core.ResultList{}
}

// MARK: Instruction Map

type InstructionMap map[string]gen.InstructionDefinition

func (m *InstructionMap) GetInstructionDefinition(
	name string,
	node *parse.InstructionNode,
) (gen.InstructionDefinition, core.ResultList) {
	inst, ok := (*m)[name]
	if !ok {
		v := (*core.UnmanagedSourceView)(nil)
		if node != nil {
			v = &node.Operator
		}

		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "undefined instruction",
			Location: v,
		}})
	}

	return inst, core.ResultList{}
}

func PrepareTestForInstructionGeneration(
	src core.SourceView,
	t *testing.T,
) (parse.InstructionNode, *gen.FunctionGenerationContext) {
	t.Helper()

	tkns, err := lex.NewTokenizer().Tokenize(src)
	tknView := parse.NewTokenView(tkns)
	assert.NoError(t, err)

	// TODO: do no use parser here? test only the instruction set unit.
	node, result := parse.NewInstructionParser().Parse(&tknView)
	assert.Nil(t, result)

	intType := gen.NewNamedTypeInfo("$32", big.NewInt(32), nil)
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

	return node, ctx
}

func TestInstructionCreateTarget(t *testing.T) {
	src := core.NewSourceView("$32 %c = add %a %b\n")
	node, ctx := PrepareTestForInstructionGeneration(src, t)

	generator := gen.NewInstructionGenerator()
	_, results := generator.Generate(ctx, node)
	assert.True(t, results.IsEmpty())

	registers := ctx.Registers
	intType := ctx.Types.GetType("$32")
	assert.NotNil(t, intType)

	target := registers.GetRegister("%c")
	assert.NotNil(t, target)
	assert.Equal(t, "%c", target.Name)
	assert.Equal(t, intType, target.Type.Base)
	assert.Equal(t, src.Unmanaged().Subview(0, 6), target.Declaration)
}

func TestUndefinedTargetType(t *testing.T) {
	src := core.NewSourceView("$undefined %c = ADD %a %b\n")
	node, ctx := PrepareTestForInstructionGeneration(src, t)

	generator := gen.NewInstructionGenerator()
	info, results := generator.Generate(ctx, node)
	assert.False(t, results.IsEmpty())
	assert.Nil(t, info)
}
