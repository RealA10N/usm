package managers

import (
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

func NewManagerCreators() gen.ManagerCreators {
	return gen.ManagerCreators{
		RegisterManagerCreator: NewRegisterManager,
		LabelManagerCreator:    NewLabelManager,
		TypeManagerCreator:     NewTypeManager,
	}
}

func NewGenerationContext() *gen.GenerationContext[usm64core.Instruction] {
	return &gen.GenerationContext[usm64core.Instruction]{
		ManagerCreators: NewManagerCreators(),
		Instructions:    NewInstructionManager(),
		PointerSize:     8,
	}
}
