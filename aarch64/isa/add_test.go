package aarch64isa_test

import (
	"testing"

	"fmt"

	"alon.kr/x/aarch64codegen/immediates"
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/aarch64codegen/registers"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64isa "alon.kr/x/usm/aarch64/isa"
	aarch64managers "alon.kr/x/usm/aarch64/managers"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func buildInstructionFromSource(
	t *testing.T,
	def gen.InstructionDefinition,
	src string,
) (*gen.InstructionInfo, aarch64codegen.Instruction) {
	srcView := core.NewSourceView(src)

	tokenizer := lex.NewTokenizer()
	tkns, err := tokenizer.Tokenize(srcView)
	assert.NoError(t, err)

	tknsView := parse.NewTokenView(tkns)
	p := parse.NewInstructionParser()
	node, result := p.Parse(&tknsView)
	assert.Nil(t, result)

	ctx := aarch64managers.NewGenerationContext().
		NewFileGenerationContext(srcView.Ctx()).
		NewFunctionGenerationContext()

	generator := gen.NewInstructionGenerator()
	info, results := generator.Generate(ctx, node)
	assert.True(t, results.IsEmpty())
	assert.NotNil(t, info)

	inst, ok := info.Definition.(aarch64codegen.Instruction)
	assert.True(t, ok)

	return info, inst
}

func assertExpectedCodegen(
	t *testing.T,
	def gen.InstructionDefinition,
	expected instructions.Instruction,
	src string,
) {
	info, inst := buildInstructionFromSource(t, def, src)

	generationContext := &aarch64codegen.InstructionCodegenContext{
		InstructionInfo: info,
	}

	code, results := inst.Codegen(generationContext)
	assert.True(t, results.IsEmpty())

	assert.Equal(t, expected.Binary(), code.Binary())
}

func TestAddExpectedCodegen(t *testing.T) {
	def := aarch64isa.NewAdd()

	testCases := []struct {
		src      string
		expected instructions.Instruction
	}{
		{
			"%x0 = add %x1 %x2\n",
			instructions.NewAddShiftedRegister(
				registers.GPRegisterX0,
				registers.GPRegisterX1,
				registers.GPRegisterX2,
			),
		},
		{
			"%xzr = add %xzr %xzr\n",
			instructions.NewAddShiftedRegister(
				registers.GPRegisterXZR,
				registers.GPRegisterXZR,
				registers.GPRegisterXZR,
			),
		},
		{
			"%x0 = add %x1 $12 #1234\n",
			instructions.NewAddImmediate(
				registers.GPorSPRegisterX0,
				registers.GPorSPRegisterX1,
				immediates.Immediate12(1234),
			),
		},
		{
			"%sp = add %sp $12 #0xfff\n",
			instructions.NewAddImmediate(
				registers.GPorSPRegisterSP,
				registers.GPorSPRegisterSP,
				immediates.Immediate12(0xfff),
			),
		},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprint(idx), func(t *testing.T) {
			assertExpectedCodegen(t, def, testCase.expected, testCase.src)
		})
	}
}
