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

type RegisterArgumentGenerator struct {
	RegisterGenerator FunctionContextGenerator[parse.RegisterNode, *RegisterInfo]
}

func NewRegisterArgumentGenerator() FunctionContextGenerator[parse.RegisterNode, ArgumentInfo] {
	return FunctionContextGenerator[parse.RegisterNode, ArgumentInfo](
		&RegisterArgumentGenerator{
			RegisterGenerator: NewRegisterGenerator(),
		},
	)
}

func (g *RegisterArgumentGenerator) Generate(
	ctx *FunctionGenerationContext,
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
