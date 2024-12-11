package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type RegisterArgumentInfo struct {
	Register *RegisterInfo
}

func (i *RegisterArgumentInfo) GetType() ReferencedTypeInfo {
	return i.Register.Type
}

type RegisterArgumentGenerator[InstT BaseInstruction] struct {
	RegisterGenerator FunctionContextGenerator[InstT, parse.RegisterNode, *RegisterInfo]
}

func NewRegisterArgumentGenerator[InstT BaseInstruction]() FunctionContextGenerator[InstT, parse.RegisterNode, ArgumentInfo] {
	return FunctionContextGenerator[InstT, parse.RegisterNode, ArgumentInfo](
		&RegisterArgumentGenerator[InstT]{
			RegisterGenerator: NewRegisterGenerator[InstT](),
		},
	)
}

func (g *RegisterArgumentGenerator[InstT]) Generate(
	ctx *FunctionGenerationContext[InstT],
	node parse.RegisterNode,
) (ArgumentInfo, core.ResultList) {
	register, results := g.RegisterGenerator.Generate(ctx, node)
	if !results.IsEmpty() {
		return nil, results
	}

	argument := RegisterArgumentInfo{
		Register: register,
	}

	return &argument, core.ResultList{}
}
