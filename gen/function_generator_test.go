package gen_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func TestSimpleFunctionGeneration(t *testing.T) {
	src := core.NewSourceView(
		`func $32 @add $32 %a {
			%b = ADD %a $32 #1
			%c = ADD %b %a
			RET
		}`,
	)
	tkns, err := lex.NewTokenizer().Tokenize(src)
	assert.NoError(t, err)

	tknView := parse.NewTokenView(tkns)
	node, result := parse.NewFunctionParser().Parse(&tknView)
	assert.Nil(t, result)

	intType := &gen.NamedTypeInfo{Name: "$32", Size: 4}

	ctx := &gen.FileGenerationContext{
		GenerationContext: &testGenerationContext,
		SourceContext:     src.Ctx(),
		Types:             &TypeMap{intType.Name: intType},
	}

	generator := gen.NewFunctionGenerator()
	function, results := generator.Generate(ctx, node)
	assert.True(t, results.IsEmpty())

	assert.NotNil(t, function.EntryBlock)
	assert.Nil(t, function.EntryBlock.NextBlock)

	registers := function.Registers
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
}

func TestIfElseFunctionGeneration(t *testing.T) {
	src := core.NewSourceView(
		`func @toBool $32 %n {
			JZ %n .zero
		.nonzero
			%bool = ADD $32 #1 $32 #0
			JMP .end
		.zero
			%bool = ADD $32 #0 $32 #0
		.end
			RET
		}`,
	)
	tkns, err := lex.NewTokenizer().Tokenize(src)
	assert.NoError(t, err)

	tknView := parse.NewTokenView(tkns)
	node, result := parse.NewFunctionParser().Parse(&tknView)
	assert.Nil(t, result)

	intType := &gen.NamedTypeInfo{Name: "$32", Size: 4}

	ctx := &gen.FileGenerationContext{
		GenerationContext: &testGenerationContext,
		SourceContext:     src.Ctx(),
		Types:             &TypeMap{intType.Name: intType},
	}

	generator := gen.NewFunctionGenerator()
	function, results := generator.Generate(ctx, node)
	assert.True(t, results.IsEmpty())

	entryBlock := function.EntryBlock
	nonzeroBlock := entryBlock.NextBlock
	zeroBlock := nonzeroBlock.NextBlock
	endBlock := zeroBlock.NextBlock

	assert.Nil(t, endBlock.NextBlock)

	assert.ElementsMatch(t, entryBlock.ForwardEdges, []*gen.BasicBlockInfo{nonzeroBlock, zeroBlock})
}
