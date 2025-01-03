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

func (i *RegisterArgumentInfo) String() string {
	return i.Register.String()
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

func NewRegisterArgumentGenerator() InstructionContextGenerator[parse.RegisterNode, ArgumentInfo] {
	return InstructionContextGenerator[parse.RegisterNode, ArgumentInfo](
		&RegisterArgumentGenerator{
			RegisterGenerator: NewRegisterGenerator(),
		},
	)
}

func (g *RegisterArgumentGenerator) Generate(
	ctx *InstructionGenerationContext,
	node parse.RegisterNode,
) (ArgumentInfo, core.ResultList) {
	register, results := g.RegisterGenerator.Generate(ctx.FunctionGenerationContext, node)
	if !results.IsEmpty() {
		return nil, results
	}

	argument := RegisterArgumentInfo{
		Register:    register,
		declaration: node.View(),
	}

	register.AddUsage(ctx.InstructionInfo)
	return &argument, core.ResultList{}
}
