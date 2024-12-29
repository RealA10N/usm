package managers

import (
	"alon.kr/x/usm/gen"
)

func NewManagerCreators() gen.ManagerCreators {
	return gen.ManagerCreators{
		RegisterManagerCreator: NewRegisterManager,
		LabelManagerCreator:    NewLabelManager,
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
