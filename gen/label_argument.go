package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

type LabelArgumentInfo[InstT BaseInstruction] struct {
	Label LabelInfo
}

func (i *LabelArgumentInfo[InstT]) GetType() ReferencedTypeInfo {
	// TODO: either (a) change implementation or (b) remove GetType() from
	// the argument interface and think of something clever.
	return ReferencedTypeInfo{}
}

// MARK: Generator

type LabelArgumentGenerator[InstT BaseInstruction] struct{}

func NewLabelArgumentGenerator[InstT BaseInstruction]() FunctionContextGenerator[InstT, parse.LabelNode, ArgumentInfo] {
	return FunctionContextGenerator[InstT, parse.LabelNode, ArgumentInfo](
		&LabelArgumentGenerator[InstT]{},
	)
}

func (g *LabelArgumentGenerator[InstT]) Generate(
	ctx *FunctionGenerationContext[InstT],
	node parse.LabelNode,
) (ArgumentInfo, core.ResultList) {
	name := nodeToSourceString(ctx.FileGenerationContext, node)
	labelInfo := ctx.Labels.GetLabel(name)

	if labelInfo == nil {
		v := node.View()
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Undefined label",
				Location: &v,
			},
		})
	}

	argument := &LabelArgumentInfo[InstT]{
		Label: *labelInfo,
	}

	return argument, core.ResultList{}
}
