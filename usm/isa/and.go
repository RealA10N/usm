package usmisa

import (
	"math/big"

	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type And struct {
	// Inherits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewAnd() gen.InstructionDefinition {
	return And{}
}

func (And) Operator(*gen.InstructionInfo) string {
	return "and"
}

func (And) PropagateConstants(info *gen.InstructionInfo) []opt.ConstantDefinition {
	return foldBinaryConstants(info, func(l, r *big.Int) *big.Int {
		return new(big.Int).And(l, r)
	})
}
