package gen_test

import (
	"math/big"
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func generateFunctionFromSource(
	t *testing.T,
	source string,
) (*gen.FunctionInfo, core.ResultList) {
	t.Helper()

	sourceView := core.NewSourceView(source)

	tkns, err := lex.NewTokenizer().Tokenize(sourceView)
	assert.NoError(t, err)

	tknView := parse.NewTokenView(tkns)
	node, result := parse.NewFunctionParser().Parse(&tknView)
	assert.Nil(t, result)

	intType := gen.NewNamedTypeInfo("$32", big.NewInt(32), nil)

	ctx := &gen.FileGenerationContext{
		GenerationContext: &testGenerationContext,
		SourceContext:     sourceView.Ctx(),
		Types:             &TypeMap{intType.Name: intType},
	}

	generator := gen.NewFunctionGenerator()
	return generator.Generate(ctx, node)
}

func TestSimpleFunctionGeneration(t *testing.T) {
	src := `func $32 @add $32 %a {
.entry
	$32 %b = ADD %a $32 #1
	$32 %c = ADD %b %a
	RET
}`

	function, results := generateFunctionFromSource(t, src)
	assert.True(t, results.IsEmpty())

	assert.NotNil(t, function.EntryBlock)

	registers := function.Registers.GetAllRegisters()
	assert.Len(t, registers, 3)

	assert.ElementsMatch(
		t,
		[][]*gen.InstructionInfo{
			nil, // TODO: make this not implementation dependent.
			{function.EntryBlock.Instructions[0]},
			{function.EntryBlock.Instructions[1]},
		},
		[][]*gen.InstructionInfo{
			registers[0].Definitions,
			registers[1].Definitions,
			registers[2].Definitions,
		},
	)

	assert.ElementsMatch(
		t,
		[][]*gen.InstructionInfo{
			{function.EntryBlock.Instructions[0], function.EntryBlock.Instructions[1]},
			{function.EntryBlock.Instructions[1]},
			nil, // TODO: make this not implementation dependent.
		},
		[][]*gen.InstructionInfo{
			registers[0].Usages,
			registers[1].Usages,
			registers[2].Usages,
		},
	)

	assert.Len(t, function.Targets, 1)
	target := function.Targets[0]
	assert.True(t, target.IsPure())
	assert.Equal(t, "$32", target.Base.Name)

	assert.Equal(t, src, function.String())
}

func TestIfElseFunctionGeneration(t *testing.T) {
	src := `func @toBool $32 %n {
.entry
	JZ %n .zero
.nonzero
	$32 %bool = ADD $32 #1 $32 #0
	JMP .end
.zero
	$32 %bool = ADD $32 #0 $32 #0
.end
	RET
}`

	function, results := generateFunctionFromSource(t, src)
	assert.True(t, results.IsEmpty())

	blocks := function.CollectBasicBlocks()
	assert.Len(t, blocks, 4)

	entryBlock := blocks[0]
	nonzeroBlock := blocks[1]
	zeroBlock := blocks[2]
	// endBlock := blocks[3]

	assert.ElementsMatch(
		t,
		entryBlock.ForwardEdges,
		[]*gen.BasicBlockInfo{nonzeroBlock, zeroBlock},
	)

	assert.Equal(t, src, function.String())
}

func TestEmptyFunctionGeneration(t *testing.T) {
	function, results := generateFunctionFromSource(t, `func @empty { }`)
	assert.False(t, results.IsEmpty())
	assert.Nil(t, function)
	details := results.Head.Value
	assert.Contains(t, details[0].Message, "at least one instruction")
}

func TestNoReturnFunctionGeneration(t *testing.T) {
	src := `func @noReturn {
				$32 %n = ADD $32 #1 $32 #2
			}`
	function, results := generateFunctionFromSource(t, src)
	assert.False(t, results.IsEmpty())
	assert.Nil(t, function)
	details := results.Head.Value
	assert.Contains(t, details[0].Message, "end a function")
}

func TestNoExplicitRegisterType(t *testing.T) {
	src := `func @noExplicitType {
				%a = ADD $32 #1 $32 #2
				RET
			}`

	function, results := generateFunctionFromSource(t, src)
	assert.False(t, results.IsEmpty())
	assert.Nil(t, function)
	details := results.Head.Value
	assert.Contains(t, details[0].Message, "untyped register")
}

func TestExplicitRegisterDefinitionNotOnSecondSight(t *testing.T) {
	src := `func @main {
				%a = ADD $32 #0 $32 #0
				$32 %a = ADD %a $32 #1
				RET
			}`
	function, results := generateFunctionFromSource(t, src)
	assert.True(t, results.IsEmpty())
	assert.NotNil(t, function)

	a := function.Registers.GetRegister("%a")
	assert.NotNil(t, a)
	assert.Equal(t, "$32", a.Type.String())
}
