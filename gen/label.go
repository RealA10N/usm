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
	// The name of the label, as it appears in the source code.
	Name string

	// The index of the instruction that the label is referencing.
	InstructionIndex core.UsmUint

	// A view of the label declaration in the source code.
	Declaration core.UnmanagedSourceView
}

// MARK: Manager

type LabelManager interface {
	GetLabel(name string) *LabelInfo
	NewLabel(info LabelInfo) core.Result
}

// MARK: Generator

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
