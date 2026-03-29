package opt

import (
	"alon.kr/x/list"
	"alon.kr/x/set"
	"alon.kr/x/stack"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/transform"
)

// An instruction is considered useful if it is critical.
// A register is considered useful if it is an argument of a useful instruction.
// If a register is useful, the instruction that defines it is also considered
// useful.
//
// From the above, we get the following algorithm:
// Construct a "useful instructions" stack.
// First, push all critical instructions into the stack.
// For each useful instruction, for all arguments to the instruction, push the
// instruction that defines the argument into the stack (It is a single, unique
// instruction since we assume SSA form).

type DCESupportedInstruction interface {
	gen.InstructionDefinition

	// Returns true if the instruction is a critical instruction, which means
	// it can't be removed by the dead code elimination process, by definition.
	//
	// A critical instruction might be a function call, a branch, an instruction
	// with a side effect, etc.
	IsCritical(info *gen.InstructionInfo) bool

	// Returns the register argument slots that the instruction writes.
	Defines(info *gen.InstructionInfo) []*gen.RegisterArgumentInfo

	// Returns the register argument slots that the instruction reads.
	Uses(info *gen.InstructionInfo) []*gen.RegisterArgumentInfo
}

func newDCENotSupportedError(instruction *gen.InstructionInfo) core.ResultList {
	return list.FromSingle(core.Result{{
		Type:     core.InternalErrorResult,
		Message:  "Instruction does not support dead code elimination",
		Location: instruction.Declaration,
	}})
}

type CriticalInstruction struct{}

func (CriticalInstruction) IsCritical(*gen.InstructionInfo) bool {
	return true
}

type NonCriticalInstruction struct{}

func (NonCriticalInstruction) IsCritical(*gen.InstructionInfo) bool {
	return false
}

type UsesInstruction struct{}

func (UsesInstruction) Uses(info *gen.InstructionInfo) []*gen.RegisterArgumentInfo {
	result := make([]*gen.RegisterArgumentInfo, 0, len(info.Arguments))
	for _, a := range info.Arguments {
		if regArg, ok := a.(*gen.RegisterArgumentInfo); ok {
			result = append(result, regArg)
		}
	}
	return result
}

type UsesNothingInstruction struct{}

func (UsesNothingInstruction) Uses(*gen.InstructionInfo) []*gen.RegisterArgumentInfo {
	return []*gen.RegisterArgumentInfo{}
}

type DefinesTargetsInstruction struct{}

func (DefinesTargetsInstruction) Defines(info *gen.InstructionInfo) []*gen.RegisterArgumentInfo {
	result := make([]*gen.RegisterArgumentInfo, 0, len(info.Targets))
	for _, t := range info.Targets {
		if regArg, ok := t.(*gen.RegisterArgumentInfo); ok {
			result = append(result, regArg)
		}
	}
	return result
}

type DefinesNothingInstruction struct{}

func (DefinesNothingInstruction) Defines(*gen.InstructionInfo) []*gen.RegisterArgumentInfo {
	return []*gen.RegisterArgumentInfo{}
}

func collectCriticalInstructions(
	function *gen.FunctionInfo,
) (stack.Stack[*gen.InstructionInfo], core.ResultList) {
	collected := stack.New[*gen.InstructionInfo]()
	results := core.ResultList{}

	for block := function.EntryBlock; block != nil; block = block.NextBlock {
		for _, instruction := range block.Instructions {
			dceInstruction, ok := instruction.Definition.(DCESupportedInstruction)
			if !ok {
				curResults := newDCENotSupportedError(instruction)
				results.Extend(&curResults)
			} else if dceInstruction.IsCritical(instruction) {
				collected.Push(instruction)
			}
		}
	}

	return collected, results
}

func collectUsefulInstructions(
	function *gen.FunctionInfo,
) (set.Set[*gen.InstructionInfo], core.ResultList) {
	usefulInstructions := set.New[*gen.InstructionInfo]()
	processedInstructions := set.New[*gen.InstructionInfo]()
	unprocessedInstructions, results := collectCriticalInstructions(function)
	if !results.IsEmpty() {
		return nil, results
	}

	for len(unprocessedInstructions) > 0 {
		instruction := unprocessedInstructions.Top()
		unprocessedInstructions.Pop()
		processedInstructions.Add(instruction)

		dceInstruction, ok := instruction.Definition.(DCESupportedInstruction)
		if !ok {
			// Should not happen, since we try to convert all instructions to
			// DCESupportedInstruction in the collection of critical instructions
			// phase.
			return nil, newDCENotSupportedError(instruction)
		}

		usefulInstructions.Add(instruction)

		for _, useArg := range dceInstruction.Uses(instruction) {
			register := useArg.Register
			// In SSA form, a register can have at most one definition, and thus
			// in SSA form this optimization can be very effective.
			// In the general sense however, we do not know what definition(s)
			// actually reach the this useful use of the register, so the best
			// we can do is to treat all references that define it as
			// potentially reaching the use.
			for _, ref := range register.References {
				refDCE, ok := ref.Definition.(DCESupportedInstruction)
				if !ok {
					return nil, newDCENotSupportedError(ref)
				}
				for _, defArg := range refDCE.Defines(ref) {
					if defArg.Register == register {
						if !processedInstructions.Contains(ref) {
							unprocessedInstructions.Push(ref)
						}
						break
					}
				}
			}
		}
	}

	return usefulInstructions, core.ResultList{}
}

func FileToDeadCodeElimination(file *gen.FileInfo) core.ResultList {
	results := core.ResultList{}

	for _, function := range file.Functions {
		if function.IsDefined() {
			curResults := DeadCodeElimination(function)
			results.Extend(&curResults)
		}
	}

	return results
}

func TransformFileToDeadCodeElimination(
	data *transform.TargetData,
) (*transform.TargetData, core.ResultList) {
	results := FileToDeadCodeElimination(data.Code)
	return data, results
}

func DeadCodeElimination(function *gen.FunctionInfo) core.ResultList {
	usefulInstructions, results := collectUsefulInstructions(function)
	if !results.IsEmpty() {
		return results
	}

	for block := function.EntryBlock; block != nil; block = block.NextBlock {
		// We remove instructions in reverse order to avoid invalidating the
		// indices of the remaining instructions.
		for i := len(block.Instructions) - 1; i >= 0; i-- {
			instruction := block.Instructions[i]
			if !usefulInstructions.Contains(instruction) {
				ok := block.RemoveInstruction(instruction)
				if !ok {
					results.Append(core.Result{{
						Type:     core.InternalErrorResult,
						Message:  "Failed to eliminate instruction",
						Location: instruction.Declaration,
					}})
				}
			}
		}
	}

	return results
}
