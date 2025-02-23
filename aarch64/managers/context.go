package aarch64managers

import "alon.kr/x/usm/gen"

func NewManagerCreators() gen.ManagerCreators {
	return gen.ManagerCreators{
		RegisterManagerCreator: NewRegisterManager,
		LabelManagerCreator:    gen.NewLabelMap,
		TypeManagerCreator:     NewTypeManager,
	}
}

func NewGenerationContext() *gen.GenerationContext {
	return &gen.GenerationContext{
		ManagerCreators: NewManagerCreators(),
		Instructions:    NewInstructionManager(),
		PointerSize:     8,
	}
}
