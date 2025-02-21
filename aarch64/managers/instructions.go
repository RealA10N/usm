package aarch64managers

import (
	"alon.kr/x/faststringmap"
	aarch64isa "alon.kr/x/usm/aarch64/isa"
	"alon.kr/x/usm/gen"
)

func NewInstructionManager() gen.InstructionManager {
	return gen.NewInstructionMap(
		[]faststringmap.MapEntry[gen.InstructionDefinition]{
			{Key: "movz", Value: aarch64isa.NewMovzInstructionDefinition()},
		},
		false,
	)
}
