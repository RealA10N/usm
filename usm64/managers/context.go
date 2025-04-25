package managers

import (
	"math/big"

	"alon.kr/x/usm/gen"
)

func NewManagerCreators() gen.ManagerCreators {
	return gen.ManagerCreators{
		RegisterManagerCreator: NewRegisterManager,
		LabelManagerCreator:    gen.NewLabelMap,
		TypeManagerCreator:     NewTypeManager,
		GlobalManagerCreator:   gen.NewGlobalMap,
	}
}

func NewGenerationContext() *gen.GenerationContext {
	return &gen.GenerationContext{
		ManagerCreators: NewManagerCreators(),
		Instructions:    NewInstructionManager(),
		PointerSize:     big.NewInt(64),
	}
}
