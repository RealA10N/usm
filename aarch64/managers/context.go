package aarch64managers

import (
	"math/big"

	"alon.kr/x/usm/gen"
	usmmanagers "alon.kr/x/usm/usm/managers"
)

func NewManagerCreators() gen.ManagerCreators {
	return gen.ManagerCreators{
		RegisterManagerCreator: NewRegisterManager,
		LabelManagerCreator:    gen.NewLabelMap,
		TypeManagerCreator:     NewTypeManager,
		GlobalManagerCreator:   gen.NewGlobalMap,
		VariableManagerCreator: usmmanagers.NewVariableManager,
	}
}

func NewGenerationContext() *gen.GenerationContext {
	return &gen.GenerationContext{
		ManagerCreators: NewManagerCreators(),
		Instructions:    NewInstructionManager(),
		PointerSize:     big.NewInt(64),
	}
}
