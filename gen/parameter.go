package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type ParameterGenerator[InstT BaseInstruction] struct {
	RegisterGenerator       Generator[InstT, parse.RegisterNode, *RegisterInfo]
	ReferencedTypeGenerator Generator[InstT, parse.TypeNode, ReferencedTypeInfo]
}

func NewParameterGenerator[InstT BaseInstruction]() Generator[InstT, parse.ParameterNode, *RegisterInfo] {
	return Generator[InstT, parse.ParameterNode, *RegisterInfo](
		&ParameterGenerator[InstT]{
			RegisterGenerator:       NewRegisterGenerator[InstT](),
			ReferencedTypeGenerator: NewReferencedTypeGenerator[InstT](),
		},
	)
}

func (g *ParameterGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.ParameterNode,
) (*RegisterInfo, core.ResultList) {
	results := core.ResultList{}

	typeInfo, typeResults := g.ReferencedTypeGenerator.Generate(ctx, node.Type)
	results.Extend(&typeResults)

	registerInfo, registerResults := g.RegisterGenerator.Generate(ctx, node.Register)
	results.Extend(&registerResults)

	if !results.IsEmpty() {
		return nil, results
	}

	if !typeInfo.Equals(registerInfo.Type) {
		v := node.View()
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Register previously declared with different type",
				Location: &v,
			},
			{
				Type:     core.HintResult,
				Message:  "Previously declared here",
				Location: &registerInfo.Declaration,
			},
		})
	}

	return registerInfo, core.ResultList{}
}
