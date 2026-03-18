package gen

import "alon.kr/x/usm/core"

// MARK: Info

type RegisterTargetInfo struct {
	Register *RegisterInfo

	declaration *core.UnmanagedSourceView
}

func NewRegisterTargetInfo(register *RegisterInfo) *RegisterTargetInfo {
	return &RegisterTargetInfo{
		Register: register,
	}
}

func (i *RegisterTargetInfo) String() string {
	return i.Register.Type.String() + " " + i.Register.String()
}

func (i *RegisterTargetInfo) Declaration() *core.UnmanagedSourceView {
	return i.declaration
}

func (i *RegisterTargetInfo) OnAttach(instruction *InstructionInfo) {
	i.Register.AddDefinition(instruction)
}

func (i *RegisterTargetInfo) OnDetach(instruction *InstructionInfo) {
	i.Register.RemoveDefinition(instruction)
}

// Switch the target to use a different register, instead of the current one,
// updating the Definitions lists on both the old and new register accordingly.
// The instruction parameter must be the InstructionInfo that owns this target.
func (i *RegisterTargetInfo) SwitchRegister(
	instruction *InstructionInfo,
	newRegister *RegisterInfo,
) {
	i.Register.RemoveDefinition(instruction)
	i.Register = newRegister
	i.Register.AddDefinition(instruction)
}
