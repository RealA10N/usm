package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type GlobalArgumentGenerator struct{}

func NewGlobalArgumentGenerator() InstructionContextGenerator[parse.GlobalNode, ArgumentInfo] {
	return InstructionContextGenerator[parse.GlobalNode, ArgumentInfo](
		&GlobalArgumentGenerator{},
	)
}

func (g *GlobalArgumentGenerator) Generate(
	ctx *InstructionGenerationContext,
	node parse.GlobalNode,
) (ArgumentInfo, core.ResultList) {
	name := NodeToSourceString(ctx.FileGenerationContext, node)
	global := ctx.Globals.GetGlobal(name)

	if global == nil {
		v := node.View()
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Undefined global",
				Location: &v,
			},
		})
	}

	argument := &GlobalArgumentInfo{
		GlobalInfo:  global,
		declaration: &node.UnmanagedSourceView,
	}

	return argument, core.ResultList{}
}
