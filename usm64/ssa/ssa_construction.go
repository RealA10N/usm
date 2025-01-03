package usm64ssa

import (
	"fmt"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/ssa"
	usm64isa "alon.kr/x/usm/usm64/isa"
)

type ConstructionScheme struct {
	RenamesPerRegister map[*gen.RegisterInfo]uint
}

func NewConstructionScheme() ssa.SsaConstructionScheme {
	scheme := &ConstructionScheme{
		RenamesPerRegister: make(map[*gen.RegisterInfo]uint),
	}

	return ssa.SsaConstructionScheme(scheme)
}

func (s *ConstructionScheme) NewPhiInstruction(
	block *gen.BasicBlockInfo,
	register *gen.RegisterInfo,
) (ssa.PhiInstruction, core.ResultList) {
	info := gen.NewEmptyInstructionInfo(nil)
	target := gen.NewRegisterArgument(register)
	info.AppendTarget(&target)
	instruction, results := usm64isa.NewPhiInstruction(info)
	info.SetBaseInstruction(instruction)
	block.PrependInstruction(info)
	return ssa.PhiInstruction(instruction), results
}

func (s *ConstructionScheme) NewRenamedRegister(
	register *gen.RegisterInfo,
) *gen.RegisterInfo {
	renamedNumber := s.RenamesPerRegister[register]
	s.RenamesPerRegister[register]++
	renamedName := register.Name + "_" + fmt.Sprint(renamedNumber)
	return gen.NewRegisterInfo(renamedName, register.Type)
}

func (s *ConstructionScheme) renameArgument(
	argument gen.ArgumentInfo,
	reachingSet ssa.ReachingDefinitionsSet,
) core.ResultList {
	if argument, ok := argument.(*gen.RegisterArgumentInfo); ok {
		baseRegister := argument.Register
		renamedRegister := reachingSet.GetReachingDefinition(baseRegister)
		argument.SwitchRegister(renamedRegister)
	}

	return core.ResultList{}
}

func (s *ConstructionScheme) renameTarget(
	target *gen.RegisterArgumentInfo,
	reachingSet ssa.ReachingDefinitionsSet,
) core.ResultList {
	baseRegister := target.Register
	renamedRegister := reachingSet.RenameDefinitionRegister(baseRegister)
	target.SwitchRegister(renamedRegister)
	return core.ResultList{}
}

func (s *ConstructionScheme) renameInstruction(
	instruction *gen.InstructionInfo,
	reachingSet ssa.ReachingDefinitionsSet,
) core.ResultList {
	// First, we rename the arguments.
	for _, argument := range instruction.Arguments {
		results := s.renameArgument(argument, reachingSet)
		if !results.IsEmpty() {
			return results
		}
	}

	// Then, we define the new registers (targets).
	for _, target := range instruction.Targets {
		results := s.renameTarget(target, reachingSet)
		if !results.IsEmpty() {
			return results
		}
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

func ConvertToSsaForm(function *gen.FunctionInfo) core.ResultList {
	constructionScheme := NewConstructionScheme()
	ssaInfo := ssa.NewFunctionSsaInfo(function, constructionScheme)
	results := ssaInfo.InsertPhiInstructions()
	if !results.IsEmpty() {
		return results
	}

	results = ssaInfo.RenameRegisters()
	if results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
