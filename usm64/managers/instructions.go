package managers

import (
	"strings"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/parse"
	usm64isa "alon.kr/x/usm/usm64/isa"
)

// TODO: optimization: implement using alon.kr/x/faststringmap
type InstructionMap map[string]gen.InstructionDefinition

func (m *InstructionMap) GetInstructionDefinition(
	name string,
	node parse.InstructionNode,
) (gen.InstructionDefinition, core.ResultList) {
	key := strings.ToLower(name)
	instDef, ok := (*m)[key]
	if !ok {
		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Undefined instruction",
			Location: &node.Operator,
		}})
	}
	return instDef, core.ResultList{}
}

func NewInstructionManager() gen.InstructionManager {
	return gen.InstructionManager(
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
			"put":  usm64isa.NewPutInstructionDefinition(),
			"term": usm64isa.NewTerminateInstructionDefinition(),
		},
	)
}
