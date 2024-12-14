package usm64isa

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type JumpInstruction struct {
	Label usm64core.Label
}

func (i *JumpInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) usm64core.EmulationError {
	ctx.JumpToLabel(i.Label)
	return nil
}

func NewJumpInstruction(
	targets []usm64core.Register,
	argument []usm64core.Argument,
) (usm64core.Instruction, core.ResultList) {

	label, ok := argument[0].(usm64core.Label)
	if !ok {
		v := argument[0].Declaration()
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected label argument",
				Location: &v,
			},
		})
	}

	return &JumpInstruction{Label: label}, core.ResultList{}
}

func NewJumpInstructionDefinition() gen.InstructionDefinition[usm64core.Instruction] {
	return &FixedInstructionDefinition{
		Targets:   0,
		Arguments: 1,
		Creator:   NewJumpInstruction,
	}
}
