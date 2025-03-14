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
type LabelDefinitionGenerator struct{}

func NewLabelDefinitionGenerator() FunctionContextGenerator[
	parse.LabelNode,
	*LabelInfo,
] {
	return FunctionContextGenerator[parse.LabelNode, *LabelInfo](
		&LabelDefinitionGenerator{},
	)
}

func (g *LabelDefinitionGenerator) Generate(
	ctx *FunctionGenerationContext,
	node parse.LabelNode,
) (*LabelInfo, core.ResultList) {
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

	newLabelInfo := &LabelInfo{
		Name:        name,
		BasicBlock:  nil, // Defined later in compilation.
		Declaration: declaration,
	}

	result := ctx.Labels.NewLabel(newLabelInfo)
	if result != nil {
		return nil, list.FromSingle(result)
	}

	return newLabelInfo, core.ResultList{}
}
