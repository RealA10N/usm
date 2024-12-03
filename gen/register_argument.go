package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type RegisterArgumentGenerator[InstT BaseInstruction] struct{}

func (g *RegisterArgumentGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.RegisterNode,
) (*ArgumentInfo, core.ResultList) {
	name := string(node.Raw(ctx.SourceContext))
	register := ctx.Registers.GetRegister(name)

	if register == nil {
		return nil, list.FromSingle(NewUndefinedRegisterResult(node.View()))
	}

	argument := ArgumentInfo{
		Type: register.Type,
	}

	return &argument, core.ResultList{}
}
