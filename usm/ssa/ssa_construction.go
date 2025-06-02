package usmssa

import (
	"fmt"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/opt/ssa"
	"alon.kr/x/usm/transform"
	usmisa "alon.kr/x/usm/usm/isa"
)

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

	target := gen.NewTargetInfo(register)
	info.AppendTarget(&target)

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
	instruction *gen.InstructionInfo,
	target *gen.TargetInfo,
	reachingSet ssa.ReachingDefinitionsSet,
) core.ResultList {
	baseRegister := target.Register
	renamedRegister := reachingSet.RenameDefinitionRegister(baseRegister)
	instruction.SwitchTarget(target, renamedRegister)
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
		results := s.renameTarget(instruction, target, reachingSet)
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
