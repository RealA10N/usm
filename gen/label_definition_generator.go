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

func NewLabelDefinitionGenerator[InstT BaseInstruction]() LabelContextGenerator[InstT, parse.LabelNode, LabelInfo] {
	return LabelContextGenerator[InstT, parse.LabelNode, LabelInfo](
		&LabelDefinitionGenerator[InstT]{},
	)
}

func (g *LabelDefinitionGenerator[InstT]) Generate(
	ctx *LabelGenerationContext[InstT],
	node parse.LabelNode,
) (LabelInfo, core.ResultList) {
	name := nodeToSourceString(ctx.FileGenerationContext, node)
	labelInfo := ctx.Labels.GetLabel(name)
	declaration := node.View()

	if labelInfo != nil {
		return LabelInfo{}, list.FromSingle(core.Result{
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

	newLabelInfo := LabelInfo{
		Name:             name,
		InstructionIndex: ctx.CurrentInstructionIndex,
		Declaration:      declaration,
	}

	result := ctx.Labels.NewLabel(newLabelInfo)
	if result != nil {
		return LabelInfo{}, list.FromSingle(result)
	}

	return newLabelInfo, core.ResultList{}
}
