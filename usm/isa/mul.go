package usmisa

import (
	"math/big"

	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type Mul struct {
	// Inherits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewMul() gen.InstructionDefinition {
	return Mul{}
}

func (Mul) Operator(*gen.InstructionInfo) string {
	return "mul"
}

func (Mul) PropagateConstants(info *gen.InstructionInfo) []opt.ConstantDefinition {
	return foldBinaryConstants(info, func(l, r *big.Int) *big.Int {
		return new(big.Int).Mul(l, r)
	})
}
