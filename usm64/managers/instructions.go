package managers

import (
	"alon.kr/x/faststringmap"
	"alon.kr/x/usm/gen"
	usm64isa "alon.kr/x/usm/usm64/isa"
)

func NewInstructionManager() gen.InstructionManager {
	return gen.NewInstructionMap(
		[]faststringmap.MapEntry[gen.InstructionDefinition]{
			// mov
			{Key: "", Value: usm64isa.NewMoveInstructionDefinition()},
			{Key: "mov", Value: usm64isa.NewMoveInstructionDefinition()},

			// arithmetic
			{Key: "add", Value: usm64isa.NewAddInstructionDefinition()},

			// control flow
			{Key: "j", Value: usm64isa.NewJumpInstructionDefinition()},
			{Key: "jz", Value: usm64isa.NewJumpZeroInstructionDefinition()},
			{Key: "jnz", Value: usm64isa.NewJumpNotZeroInstructionDefinition()},

			// SSA
			{Key: "phi", Value: usm64isa.NewPhiInstructionDefinition()},

			// debug
			{Key: "put", Value: usm64isa.NewPutInstructionDefinition()},
			{Key: "term", Value: usm64isa.NewTerminateInstructionDefinition()},
		},
		false,
	)
}
