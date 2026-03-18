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
	return i.Register.Type.String() + " " + i.Register.String()
}

func (i *RegisterArgumentInfo) Declaration() *core.UnmanagedSourceView {
	return i.declaration
}

func (i *RegisterArgumentInfo) OnAttach(instruction *InstructionInfo) {
	i.Register.AddUsage(instruction)
}

func (i *RegisterArgumentInfo) OnDetach(instruction *InstructionInfo) {
	i.Register.RemoveUsage(instruction)
}

// Switch the argument to use a different register, instead of the current one,
// updating the Usages lists on both the old and new register accordingly.
// The instruction parameter must be the InstructionInfo that owns this argument.
func (i *RegisterArgumentInfo) SwitchRegister(
	instruction *InstructionInfo,
	newRegister *RegisterInfo,
) {
	i.Register.RemoveUsage(instruction)
	i.Register = newRegister
	i.Register.AddUsage(instruction)
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
