package opt

import (
	"alon.kr/x/list"
	"alon.kr/x/set"
	"alon.kr/x/stack"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
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

	// Returns the registers that the instruction defines, directly or indirectly.
	Defines(info *gen.InstructionInfo) []*gen.RegisterInfo

	// Returns the registers that the instruction uses, directly or indirectly.
	Uses(info *gen.InstructionInfo) []*gen.RegisterInfo
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

type UsesArgumentsInstruction struct{}

func (UsesArgumentsInstruction) Uses(info *gen.InstructionInfo) []*gen.RegisterInfo {
	return gen.ArgumentsToRegisters(info.Arguments)
}

type UsesNothingInstruction struct{}

func (UsesNothingInstruction) Uses(*gen.InstructionInfo) []*gen.RegisterInfo {
	return []*gen.RegisterInfo{}
}

type DefinesTargetsInstruction struct{}

func (DefinesTargetsInstruction) Defines(info *gen.InstructionInfo) []*gen.RegisterInfo {
	return gen.TargetsToRegisters(info.Targets)
}

type DefinesNothingInstruction struct{}

func (DefinesNothingInstruction) Defines(*gen.InstructionInfo) []*gen.RegisterInfo {
	return []*gen.RegisterInfo{}
}

func collectCriticalInstructions(
	function *gen.FunctionInfo,
) (stack.Stack[*gen.InstructionInfo], core.ResultList) {
	collected := stack.New[*gen.InstructionInfo]()
	results := core.ResultList{}

	for block := function.EntryBlock; block != nil; block = block.NextBlock {
		for _, instruction := range block.Instructions {
			dceInstruction, ok := instruction.Instruction.(DCESupportedInstruction)
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

		dceInstruction, ok := instruction.Instruction.(DCESupportedInstruction)
		if !ok {
			// Should not happen, since we try to convert all instructions to
			// DCESupportedInstruction in the collection of critical instructions
			// phase.
			return nil, newDCENotSupportedError(instruction)
		}

		usefulInstructions.Add(instruction)

		for _, register := range dceInstruction.Uses(instruction) {
			definitions := register.Definitions

			// In SSA form, a register can have at most one definition, and thus
			// in SSA form this optimization can be very effective.
			// In the general sense however, we do not know what definition(s)
			// actually reach the this useful use of the register, so the best
			// we can do is to treat all definitions as potentially reaching the
			// use.
			for _, definition := range definitions {
				if !processedInstructions.Contains(definition) {
					unprocessedInstructions.Push(definition)
				}
			}
		}
	}

	return usefulInstructions, core.ResultList{}
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
