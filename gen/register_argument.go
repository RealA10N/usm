package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type RegisterArgumentInfo struct {
	Type ReferencedTypeInfo
}

func (i *RegisterArgumentInfo) GetType() ReferencedTypeInfo {
	return i.Type
}

type RegisterArgumentGenerator[InstT BaseInstruction] struct{}

func NewRegisterArgumentGenerator[InstT BaseInstruction]() Generator[InstT, parse.RegisterNode, ArgumentInfo] {
	return Generator[InstT, parse.RegisterNode, ArgumentInfo](
		&RegisterArgumentGenerator[InstT]{},
	)
}

func (g *RegisterArgumentGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.RegisterNode,
) (ArgumentInfo, core.ResultList) {
	name := string(node.Raw(ctx.SourceContext))
	register := ctx.Registers.GetRegister(name)

	if register == nil {
		return nil, list.FromSingle(NewUndefinedRegisterResult(node.View()))
	}

	argument := RegisterArgumentInfo{
		Type: register.Type,
	}

	return &argument, core.ResultList{}
}
