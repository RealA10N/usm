package usmisa

import (
	"math/big"

	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type Or struct {
	// Inherits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewOr() gen.InstructionDefinition {
	return Or{}
}

func (Or) Operator(*gen.InstructionInfo) string {
	return "or"
}

func (Or) PropagateConstants(info *gen.InstructionInfo) []opt.ConstantDefinition {
	return foldBinaryConstants(info, func(l, r *big.Int) *big.Int {
		return new(big.Int).Or(l, r)
	})
}
