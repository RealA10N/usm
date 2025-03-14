package usm64isa

import (
	"alon.kr/x/usm/gen"
)

type baseInstruction struct {
	// A pointer the the internal USM representation of the instruction.
	// This in turn has the representation of the arguments, targets, and types.
	*gen.InstructionInfo
}

func newBaseInstruction(info *gen.InstructionInfo) baseInstruction {
	return baseInstruction{info}
}

func (i *baseInstruction) IsCritical() bool {
	return false
}

func (i *baseInstruction) Uses() []*gen.RegisterInfo {
	arguments := i.InstructionInfo.Arguments
	registers := []*gen.RegisterInfo{}
	for _, argument := range arguments {
		if argument, ok := argument.(*gen.RegisterArgumentInfo); ok {
			registers = append(registers, argument.Register)
		}
	}

	return registers
}

func (i *baseInstruction) Defines() []*gen.RegisterInfo {
	targets := i.InstructionInfo.Targets
	registers := []*gen.RegisterInfo{}
	for _, target := range targets {
		registers = append(registers, target.Register)
	}

	return registers
}
