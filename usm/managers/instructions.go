package usmmanagers

import (
	"alon.kr/x/faststringmap"
	"alon.kr/x/usm/gen"
	usmisa "alon.kr/x/usm/usm/isa"
)

func NewInstructionManager() gen.InstructionManager {
	return gen.NewInstructionMap(
		[]faststringmap.MapEntry[gen.InstructionDefinition]{
			{Key: "", Value: usmisa.NewMove()},

			// Control Flow
			{Key: "ret", Value: usmisa.NewRet()},
			{Key: "call", Value: usmisa.NewCall()},
			{Key: "jz", Value: usmisa.NewJz()},

			// Static Single Assignment (SSA)
			{Key: "phi", Value: usmisa.NewPhi()},
		},
		false,
	)
}
