package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

type RegisterArgumentInfo struct {
	Register    *RegisterInfo
	declaration core.UnmanagedSourceView
}

func (i *RegisterArgumentInfo) GetType() *ReferencedTypeInfo {
	return &i.Register.Type
}

func (i *RegisterArgumentInfo) Declaration() core.UnmanagedSourceView {
	return i.declaration
}

// MARK: Generator

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
		Register:    register,
		declaration: node.View(),
	}

	return &argument, core.ResultList{}
}
