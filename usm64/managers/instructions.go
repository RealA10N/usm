package managers

import (
	"strings"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
	usm64isa "alon.kr/x/usm/usm64/isa"
)

// TODO: optimization: implement using alon.kr/x/faststringmap
type InstructionMap map[string]gen.InstructionDefinition[usm64core.Instruction]

func (m *InstructionMap) GetInstructionDefinition(
	name string,
) (gen.InstructionDefinition[usm64core.Instruction], core.ResultList) {
	key := strings.ToLower(name)
	instDef, ok := (*m)[key]
	if !ok {
		return nil, list.FromSingle(core.Result{{
			Type:    core.ErrorResult,
			Message: "Undefined instruction",
			// TODO: add location
		}})
	}
	return instDef, core.ResultList{}
}

func NewInstructionManager() gen.InstructionManager[usm64core.Instruction] {
	return gen.InstructionManager[usm64core.Instruction](
		&InstructionMap{

			// mov
			"":    usm64isa.NewMoveInstructionDefinition(),
			"mov": usm64isa.NewMoveInstructionDefinition(),

			// arithmetic
			"add": usm64isa.NewAddInstructionDefinition(),

			// control flow
			"j":   usm64isa.NewJumpInstructionDefinition(),
			"jz":  usm64isa.NewJumpZeroInstructionDefinition(),
			"jnz": usm64isa.NewJumpNotZeroInstructionDefinition(),

			// debug
			"put": usm64isa.NewPutInstructionDefinition(),
		},
	)
}
