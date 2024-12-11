// Labels are abstracted from the "backend".
// The "frontend" (gen module) iterates over all local function labels,
// and provides the labels interface (arguments to instructions) as pointers
// to other instructions in the same function scope.

package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

type LabelInfo struct {
	Name        string
	Declaration core.UnmanagedSourceView
}

// MARK: Manager

type LabelManager interface {
	GetLabel(name string) *LabelInfo
	NewLabel(info LabelInfo) core.Result
}

// MARK: Definition

type LabelDefinitionGenerator[InstT BaseInstruction] struct{}

func (g *LabelDefinitionGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.LabelNode,
) (LabelInfo, core.ResultList) {
	name := nodeToSourceString(ctx, node)
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
		Name:        name,
		Declaration: declaration,
	}

	return newLabelInfo, core.ResultList{}
}

// MARK: Reference

type LabelReferenceGenerator[InstT BaseInstruction] struct{}

func (g *LabelReferenceGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.LabelNode,
) (LabelInfo, core.ResultList) {
	name := nodeToSourceString(ctx, node)
	labelInfo := ctx.Labels.GetLabel(name)

	if labelInfo == nil {
		v := node.View()
		return LabelInfo{}, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Undefined label",
				Location: &v,
			},
		})
	}

	return *labelInfo, core.ResultList{}
}
