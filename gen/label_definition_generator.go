package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// Generator for label *definition* nodes.
//
// Will create and add the label to the function context,
// and return an error if a label with the same name already exists.
type LabelDefinitionGenerator[InstT BaseInstruction] struct{}

func NewLabelDefinitionGenerator[InstT BaseInstruction]() LabelContextGenerator[InstT, parse.LabelNode, *LabelInfo[InstT]] {
	return LabelContextGenerator[InstT, parse.LabelNode, *LabelInfo[InstT]](
		&LabelDefinitionGenerator[InstT]{},
	)
}

func (g *LabelDefinitionGenerator[InstT]) Generate(
	ctx *LabelGenerationContext[InstT],
	node parse.LabelNode,
) (*LabelInfo[InstT], core.ResultList) {
	name := nodeToSourceString(ctx.FileGenerationContext, node)
	labelInfo := ctx.Labels.GetLabel(name)
	declaration := node.View()

	if labelInfo != nil {
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Label already defined",
				Location: &declaration,
			},
			{
				Type:     core.HintResult,
				Message:  "Previous definition here",
				Location: &labelInfo.Declaration,
			},
		})
	}

	newLabelInfo := &LabelInfo[InstT]{
		Name:        name,
		Declaration: declaration,
	}

	result := ctx.Labels.NewLabel(newLabelInfo)
	if result != nil {
		return nil, list.FromSingle(result)
	}

	return newLabelInfo, core.ResultList{}
}
