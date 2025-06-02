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
			{Key: "movz", Value: aarch64isa.NewMovz()},

			// Arithmetic
			{Key: "add", Value: aarch64isa.NewAdd()},
			{Key: "adds", Value: aarch64isa.NewAdds()},
			{Key: "sub", Value: aarch64isa.NewSub()},
			{Key: "subs", Value: aarch64isa.NewSubs()},

			// Control flow
			{Key: "b", Value: aarch64isa.NewBranch()},
			{Key: "bl", Value: aarch64isa.NewBl()},
			{Key: "ret", Value: aarch64isa.NewRet()},

			// Conditional branches
			{Key: "b.eq", Value: aarch64isa.NewBcond(immediates.ConditionEq)},
			{Key: "b.ne", Value: aarch64isa.NewBcond(immediates.ConditionNe)},
			{Key: "b.cs", Value: aarch64isa.NewBcond(immediates.ConditionCs)},
			{Key: "b.cc", Value: aarch64isa.NewBcond(immediates.ConditionCc)},
			{Key: "b.mi", Value: aarch64isa.NewBcond(immediates.ConditionMi)},
			{Key: "b.pl", Value: aarch64isa.NewBcond(immediates.ConditionPl)},
			{Key: "b.vs", Value: aarch64isa.NewBcond(immediates.ConditionVs)},
			{Key: "b.vc", Value: aarch64isa.NewBcond(immediates.ConditionVc)},
			{Key: "b.hi", Value: aarch64isa.NewBcond(immediates.ConditionHi)},
			{Key: "b.ls", Value: aarch64isa.NewBcond(immediates.ConditionLs)},
			{Key: "b.ge", Value: aarch64isa.NewBcond(immediates.ConditionGe)},
			{Key: "b.lt", Value: aarch64isa.NewBcond(immediates.ConditionLt)},
			{Key: "b.gt", Value: aarch64isa.NewBcond(immediates.ConditionGt)},
			{Key: "b.le", Value: aarch64isa.NewBcond(immediates.ConditionLe)},
			{Key: "b.al", Value: aarch64isa.NewBcond(immediates.ConditionAl)},
			{Key: "b.nv", Value: aarch64isa.NewBcond(immediates.ConditionNv)},
		},
		false,
	)
}
