package aarch64managers

import (
	"alon.kr/x/aarch64codegen/immediates"
	"alon.kr/x/faststringmap"
	aarch64isa "alon.kr/x/usm/aarch64/isa"
	"alon.kr/x/usm/gen"
)

func NewInstructionManager() gen.InstructionManager {
	return gen.NewInstructionMap(
		[]faststringmap.MapEntry[gen.InstructionDefinition]{
			// Move
			{Key: "movz", Value: aarch64isa.NewMovzInstructionDefinition()},

			// Arithmetic
			{Key: "add", Value: aarch64isa.NewAddInstructionDefinition()},
			{Key: "adds", Value: aarch64isa.NewAddsInstructionDefinition()},
			{Key: "sub", Value: aarch64isa.NewSubInstructionDefinition()},
			{Key: "subs", Value: aarch64isa.NewSubsInstructionDefinition()},

			// Control flow
			{Key: "b", Value: aarch64isa.NewBranchInstructionDefinition()},
			{Key: "bl", Value: aarch64isa.NewBlInstructionDefinition()},
			{Key: "ret", Value: aarch64isa.NewRetInstructionDefinition()},

			// Conditional branches
			{Key: "b.eq", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionEq)},
			{Key: "b.ne", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionNe)},
			{Key: "b.cs", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionCs)},
			{Key: "b.cc", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionCc)},
			{Key: "b.mi", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionMi)},
			{Key: "b.pl", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionPl)},
			{Key: "b.vs", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionVs)},
			{Key: "b.vc", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionVc)},
			{Key: "b.hi", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionHi)},
			{Key: "b.ls", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionLs)},
			{Key: "b.ge", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionGe)},
			{Key: "b.lt", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionLt)},
			{Key: "b.gt", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionGt)},
			{Key: "b.le", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionLe)},
			{Key: "b.al", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionAl)},
			{Key: "b.nv", Value: aarch64isa.NewBcondInstructionDefinition(immediates.ConditionNv)},
		},
		false,
	)
}
