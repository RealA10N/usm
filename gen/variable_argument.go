package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

type VariableArgumentInfo struct {
	Variable    *VariableInfo
	declaration *core.UnmanagedSourceView
}

func (*VariableArgumentInfo) OnAttach(*InstructionInfo) {}
func (*VariableArgumentInfo) OnDetach(*InstructionInfo) {}

func (i *VariableArgumentInfo) Declaration() *core.UnmanagedSourceView {
	return i.declaration
}

func (i *VariableArgumentInfo) String() string {
	return i.Variable.Name
}

// MARK: Generator

type VariableArgumentGenerator struct{}

func NewVariableArgumentGenerator() InstructionContextGenerator[parse.VariableNode, ArgumentInfo] {
	return InstructionContextGenerator[parse.VariableNode, ArgumentInfo](
		&VariableArgumentGenerator{},
	)
}

func (g *VariableArgumentGenerator) Generate(
	ctx *InstructionGenerationContext,
	node parse.VariableNode,
) (ArgumentInfo, core.ResultList) {
	name := NodeToSourceString(ctx.FileGenerationContext, node)
	variable := ctx.Variables.GetVariable(name)

	if variable == nil {
		v := node.View()
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Undefined variable",
				Location: &v,
			},
		})
	}

	v := node.View()
	return &VariableArgumentInfo{
		Variable:    variable,
		declaration: &v,
	}, core.ResultList{}
}
