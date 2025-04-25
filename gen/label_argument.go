package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

type LabelArgumentInfo struct {
	Label *LabelInfo

	declaration *core.UnmanagedSourceView
}

func NewLabelArgumentInfo(
	label *LabelInfo,
) *LabelArgumentInfo {
	return &LabelArgumentInfo{
		Label:       label,
		declaration: nil,
	}
}

func (i *LabelArgumentInfo) GetType() *ReferencedTypeInfo {
	return nil // Label argument does not have a type
}

func (i *LabelArgumentInfo) Declaration() *core.UnmanagedSourceView {
	return i.declaration
}

func (i *LabelArgumentInfo) String() string {
	return i.Label.Name
}

// MARK: Generator

type LabelArgumentGenerator struct{}

func NewLabelArgumentGenerator() InstructionContextGenerator[parse.LabelNode, ArgumentInfo] {
	return InstructionContextGenerator[parse.LabelNode, ArgumentInfo](
		&LabelArgumentGenerator{},
	)
}

func (g *LabelArgumentGenerator) Generate(
	ctx *InstructionGenerationContext,
	node parse.LabelNode,
) (ArgumentInfo, core.ResultList) {
	name := NodeToSourceString(ctx.FileGenerationContext, node)
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

	v := node.View()
	argument := &LabelArgumentInfo{
		Label:       labelInfo,
		declaration: &v,
	}

	return argument, core.ResultList{}
}
