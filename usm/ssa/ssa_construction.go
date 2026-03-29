package usmssa

import (
	"fmt"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt/ssa"
	"alon.kr/x/usm/transform"
	usmisa "alon.kr/x/usm/usm/isa"
)

func newSSANotSupportedError(instruction *gen.InstructionInfo) core.ResultList {
	return list.FromSingle(core.Result{{
		Type:     core.InternalErrorResult,
		Message:  "Instruction does not support SSA construction",
		Location: instruction.Declaration,
	}})
}

type ConstructionScheme struct {
	RenamesPerRegister map[*gen.RegisterInfo]uint
}

func NewConstructionScheme() ssa.SsaConstructionScheme {
	return &ConstructionScheme{
		RenamesPerRegister: make(map[*gen.RegisterInfo]uint),
	}
}

func (s *ConstructionScheme) NewPhiInstruction(
	block *gen.BasicBlockInfo,
	register *gen.RegisterInfo,
) (*gen.InstructionInfo, core.ResultList) {
	info := gen.NewEmptyInstructionInfo(nil)
	info.SetInstruction(usmisa.NewPhi())

	target := gen.NewRegisterArgumentInfo(register)
	info.AppendTarget(target)

	block.PrependInstruction(info)

	return info, core.ResultList{}
}

func (s *ConstructionScheme) NewRenamedRegister(
	register *gen.RegisterInfo,
) *gen.RegisterInfo {
	renamedNumber := s.RenamesPerRegister[register]
	s.RenamesPerRegister[register]++
	renamedName := register.Name + "_" + fmt.Sprint(renamedNumber)
	return gen.NewRegisterInfo(renamedName, register.Type)
}

func (s *ConstructionScheme) renameInstruction(
	instruction *gen.InstructionInfo,
	reachingSet ssa.ReachingDefinitionsSet,
) core.ResultList {
	ssaInstr, ok := instruction.Definition.(ssa.SSASupportedInstruction)
	if !ok {
		return newSSANotSupportedError(instruction)
	}

	// Rename uses first: each used register is replaced with the current
	// reaching renamed definition.
	for _, useArg := range ssaInstr.Uses(instruction) {
		renamedRegister := reachingSet.GetReachingDefinition(useArg.Register)
		useArg.SwitchRegister(instruction, renamedRegister)
	}

	// Then rename definitions: each defined register gets a fresh name and
	// becomes the new reaching definition for downstream uses.
	for _, defArg := range ssaInstr.Defines(instruction) {
		renamedRegister := reachingSet.RenameDefinitionRegister(defArg.Register)
		defArg.SwitchRegister(instruction, renamedRegister)
	}

	return core.ResultList{}
}

func (s *ConstructionScheme) RenameBasicBlock(
	block *gen.BasicBlockInfo,
	reachingSet ssa.ReachingDefinitionsSet,
) core.ResultList {
	for _, instruction := range block.Instructions {
		results := s.renameInstruction(instruction, reachingSet)
		if !results.IsEmpty() {
			return results
		}
	}

	return core.ResultList{}
}

func FunctionToSsaForm(function *gen.FunctionInfo) core.ResultList {
	constructionScheme := NewConstructionScheme()
	ssaInfo := ssa.NewFunctionSsaInfo(function, constructionScheme)
	results := ssaInfo.InsertPhiInstructions()
	if !results.IsEmpty() {
		return results
	}

	results = ssaInfo.RenameRegisters()
	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}

func FileToSsaForm(file *gen.FileInfo) core.ResultList {
	results := core.ResultList{}

	for _, function := range file.Functions {
		if function.IsDefined() {
			curResults := FunctionToSsaForm(function)
			results.Extend(&curResults)
		}
	}

	return results
}

func TransformFileToSsaForm(
	data *transform.TargetData,
) (*transform.TargetData, core.ResultList) {
	results := FileToSsaForm(data.Code)
	return data, results
}
