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
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Undefined register used as argument",
			Location: &v,
		}})
	}

	argument := ArgumentInfo{
		Type: register.Type,
	}

	return &argument, core.ResultList{}
}
