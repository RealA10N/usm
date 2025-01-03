package usm64isa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type PhiInstruction struct {
	nonBranchingInstruction
}

func (i *PhiInstruction) Emulate(
	ctx *usm64core.EmulationContext,
) core.ResultList {
	labelArgument := i.InstructionInfo.Arguments[0].(*gen.LabelArgumentInfo)
	return ctx.JumpToLabel(labelArgument.Label)
}

func NewPhiInstruction(
	info *gen.InstructionInfo,
) (*PhiInstruction, core.ResultList) {
	return &PhiInstruction{newNonBranchingInstruction(info)}, core.ResultList{}
}

func (i *PhiInstruction) String() string {
	return "PHI"
}

func (i *PhiInstruction) AddForwardingRegister(
	block *gen.BasicBlockInfo,
	register *gen.RegisterInfo,
) core.ResultList {
	labelArgument := gen.NewLabelArgumentInfo(block.GetRepresentingLabel())
	registerArgument := gen.NewRegisterArgument(register)
	i.Arguments = append(i.Arguments, labelArgument, &registerArgument)
	return core.ResultList{}
}

type PhiInstructionDefinition struct{}

func (d *PhiInstructionDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	// TODO: argument validation
	return NewPhiInstruction(info)
}

func (d *PhiInstructionDefinition) InferTargetTypes(
	ctx *gen.FunctionGenerationContext,
	targets []*gen.ReferencedTypeInfo,
	arguments []*gen.ReferencedTypeInfo,
) ([]gen.ReferencedTypeInfo, core.ResultList) {
	return []gen.ReferencedTypeInfo{}, core.ResultList{}
}

func NewPhiInstructionDefinition() gen.InstructionDefinition {
	return &PhiInstructionDefinition{}
}
