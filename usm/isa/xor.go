package usmisa

import (
	"math/big"

	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type Xor struct {
	// Inherits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewXor() gen.InstructionDefinition {
	return Xor{}
}

func (Xor) Operator(*gen.InstructionInfo) string {
	return "xor"
}

func (Xor) PropagateConstants(info *gen.InstructionInfo) []opt.ConstantDefinition {
	return foldBinaryConstants(info, func(l, r *big.Int) *big.Int {
		return new(big.Int).Xor(l, r)
	})
}
