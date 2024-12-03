package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type ImmediateArgumentGenerator[InstT BaseInstruction] struct{}

func (g *ImmediateArgumentGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.ImmediateNode,
) (*ArgumentInfo, core.ResultList) {
	typeName := string(node.Type.Raw(ctx.SourceContext))
	typ := ctx.Types.GetType(typeName)
	if typ == nil {
		return nil, list.FromSingle(NewUndefinedTypeResult(node.View()))
	}

	argument := ArgumentInfo{
		Type: typ,
	}

	return &argument, core.ResultList{}
}
