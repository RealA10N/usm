package usmisa

import (
	"math/big"

	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt"
)

type Sub struct {
	// Inherits most of the functionality from BinaryCalculation
	BinaryCalculation
}

func NewSub() gen.InstructionDefinition {
	return Sub{}
}

func (Sub) Operator(*gen.InstructionInfo) string {
	return "sub"
}

func (Sub) PropagateConstants(info *gen.InstructionInfo) []opt.ConstantDefinition {
	return foldBinaryConstants(info, func(l, r *big.Int) *big.Int {
		return new(big.Int).Sub(l, r)
	})
}
