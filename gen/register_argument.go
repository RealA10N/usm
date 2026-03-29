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
	i.Register.AddReference(instruction)
}

func (i *RegisterArgumentInfo) OnDetach(instruction *InstructionInfo) {
	i.Register.RemoveReference(instruction)
}

// Switch the argument to use a different register, instead of the current one,
// updating the References lists on both the old and new register accordingly.
// The instruction parameter must be the InstructionInfo that owns this argument.
func (i *RegisterArgumentInfo) SwitchRegister(
	instruction *InstructionInfo,
	newRegister *RegisterInfo,
) {
	i.Register.RemoveReference(instruction)
	i.Register = newRegister
	i.Register.AddReference(instruction)
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

	return &argument, core.ResultList{}
}
