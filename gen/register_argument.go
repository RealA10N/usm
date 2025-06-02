package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

type RegisterArgumentInfo struct {
	Register *RegisterInfo

	declaration *core.UnmanagedSourceView
}

func NewRegisterArgumentInfo(register *RegisterInfo) *RegisterArgumentInfo {
	return &RegisterArgumentInfo{
		Register: register,
	}
}

func (i *RegisterArgumentInfo) String() string {
	return i.Register.String()
}

func (i *RegisterArgumentInfo) Declaration() *core.UnmanagedSourceView {
	return i.declaration
}

// Switch the argument to use a different register, instead of the current
// one.
func (i *RegisterArgumentInfo) SwitchRegister(newRegister *RegisterInfo) {
	// TODO: handle definitions and usages
	i.Register = newRegister
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

	v := node.View()
	argument := RegisterArgumentInfo{
		Register:    register,
		declaration: &v,
	}

	register.AddUsage(ctx.InstructionInfo)
	return &argument, core.ResultList{}
}
