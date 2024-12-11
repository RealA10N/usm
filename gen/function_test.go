package gen_test

import (
	"testing"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/stretchr/testify/assert"
)

func TestFunctionGeneration(t *testing.T) {
	src := core.NewSourceView(
		`func $32 @add $32 %a {
			%b = ADD %a $32 #1
			%c = ADD %b %a
		}`,
	)
	tkns, err := lex.NewTokenizer().Tokenize(src)
	assert.NoError(t, err)

	// TODO: do no use parser here? test only the instruction set unit.
	tknView := parse.NewTokenView(tkns)
	node, result := parse.NewFunctionParser().Parse(&tknView)
	assert.Nil(t, result)

	instructions := InstructionMap{
		"ADD": &AddInstructionDefinition{},
	}

	intType := &gen.NamedTypeInfo{Name: "$32", Size: 4}

	ctx := &gen.FunctionGenerationContext[Instruction]{
		FileGenerationContext: &gen.FileGenerationContext{
			GenerationContext: &gen.GenerationContext{
				PointerSize: 8,
			},
			SourceContext: src.Ctx(),
			Types:         &TypeMap{intType.Name: intType},
		},
		Registers:    &RegisterMap{},
		Instructions: &instructions,
	}

	generator := gen.NewFunctionGenerator[Instruction]()
	_, results := generator.Generate(ctx, node)
	assert.True(t, results.IsEmpty())
}
