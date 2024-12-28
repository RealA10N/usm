package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type ArgumentGenerator struct {
	RegisterArgumentGenerator  FunctionContextGenerator[parse.RegisterNode, ArgumentInfo]
	ImmediateArgumentGenerator FunctionContextGenerator[parse.ImmediateNode, ArgumentInfo]
	LabelArgumentGenerator     FunctionContextGenerator[parse.LabelNode, ArgumentInfo]
}

func NewArgumentGenerator() FunctionContextGenerator[parse.ArgumentNode, ArgumentInfo] {
	return FunctionContextGenerator[parse.ArgumentNode, ArgumentInfo](
		&ArgumentGenerator{
			RegisterArgumentGenerator:  NewRegisterArgumentGenerator(),
			ImmediateArgumentGenerator: NewImmediateArgumentGenerator(),
			LabelArgumentGenerator:     NewLabelArgumentGenerator(),
		},
	)
}

func (g *ArgumentGenerator) Generate(
	ctx *FunctionGenerationContext,
	node parse.ArgumentNode,
) (ArgumentInfo, core.ResultList) {
	switch typedNode := node.(type) {
	case parse.RegisterNode:
		return g.RegisterArgumentGenerator.Generate(ctx, typedNode)
	case parse.ImmediateNode:
		return g.ImmediateArgumentGenerator.Generate(ctx, typedNode)
	case parse.LabelNode:
		return g.LabelArgumentGenerator.Generate(ctx, typedNode)
	default:
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Unsupported argument type",
			Location: &v,
		}})
	}
}
