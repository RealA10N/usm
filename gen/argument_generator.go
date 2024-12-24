package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type ArgumentGenerator[InstT BaseInstruction] struct {
	RegisterArgumentGenerator  FunctionContextGenerator[InstT, parse.RegisterNode, ArgumentInfo]
	ImmediateArgumentGenerator FunctionContextGenerator[InstT, parse.ImmediateNode, ArgumentInfo]
	LabelArgumentGenerator     FunctionContextGenerator[InstT, parse.LabelNode, ArgumentInfo]
}

func NewArgumentGenerator[InstT BaseInstruction]() FunctionContextGenerator[InstT, parse.ArgumentNode, ArgumentInfo] {
	return FunctionContextGenerator[InstT, parse.ArgumentNode, ArgumentInfo](
		&ArgumentGenerator[InstT]{
			RegisterArgumentGenerator:  NewRegisterArgumentGenerator[InstT](),
			ImmediateArgumentGenerator: NewImmediateArgumentGenerator[InstT](),
			LabelArgumentGenerator:     NewLabelArgumentGenerator[InstT](),
		},
	)
}

func (g *ArgumentGenerator[InstT]) Generate(
	ctx *FunctionGenerationContext[InstT],
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