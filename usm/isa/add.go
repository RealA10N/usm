package usmisa

import (
	"math/big"

	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type Add struct {
	// Inherits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewAdd() gen.InstructionDefinition {
	return Add{}
}

func (Add) Operator(*gen.InstructionInfo) string {
	return "add"
}

func (Add) PropagateConstants(info *gen.InstructionInfo) []opt.ConstantDefinition {
	return foldBinaryConstants(info, func(l, r *big.Int) *big.Int {
		return new(big.Int).Add(l, r)
	})
}
