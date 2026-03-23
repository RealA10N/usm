package usmssa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	ssaopt "alon.kr/x/usm/opt/ssa"
	"alon.kr/x/usm/transform"
)

// collectPhiInstructions returns the phi instructions at the start of a block.
// Phi instructions are always at the beginning of a block in SSA form.
func collectPhiInstructions(block *gen.BasicBlockInfo) []*gen.InstructionInfo {
	var phis []*gen.InstructionInfo
	for _, instr := range block.Instructions {
		if _, isPhi := instr.Definition.(ssaopt.PhiInstructionDefinition); !isPhi {
			break
		}
		phis = append(phis, instr)
	}
	return phis
}

// collectCopiesForPredecessor returns the parallel copy set for all phis in
// block that have an incoming value from pred.
func collectCopiesForPredecessor(
	phis []*gen.InstructionInfo,
	pred *gen.BasicBlockInfo,
) []parallelCopyEntry {
	var copies []parallelCopyEntry
	for _, phi := range phis {
		dst := phi.Targets[0].Register
		// Phi arguments are pairs: [label0, value0, label1, value1, ...]
		for i := 0; i+1 < len(phi.Arguments); i += 2 {
			labelArg, ok := phi.Arguments[i].(*gen.LabelArgumentInfo)
			if !ok {
				continue
			}
			if labelArg.Label.BasicBlock != pred {
				continue
			}
			copies = append(copies, parallelCopyEntry{dst: dst, srcArg: phi.Arguments[i+1]})
			break
		}
	}
	return copies
}

// FunctionOutOfSsaForm removes all phi instructions from function, replacing
// them with copy instructions inserted in each predecessor block.
func FunctionOutOfSsaForm(function *gen.FunctionInfo) core.ResultList {
	splitCounter := 0
	tmpCounter := 0

	splitCriticalEdgesInFunction(function, &splitCounter)

	for block := function.EntryBlock; block != nil; block = block.NextBlock {
		phis := collectPhiInstructions(block)
		if len(phis) == 0 {
			continue
		}

		// For each predecessor, build and insert the parallel copy set.
		for _, pred := range block.BackwardEdges {
			copies := collectCopiesForPredecessor(phis, pred)
			instrs := sequentializeParallelCopies(copies, function, &tmpCounter)
			for _, instr := range instrs {
				pred.InsertBeforeTerminator(instr)
			}
		}

		// Remove the phi instructions from this block.
		for _, phi := range phis {
			for _, target := range phi.Targets {
				target.Register.RemoveDefinition(phi)
			}
			block.RemoveInstruction(phi)
		}
	}

	return core.ResultList{}
}

// FileOutOfSsaForm applies FunctionOutOfSsaForm to every defined function.
func FileOutOfSsaForm(file *gen.FileInfo) core.ResultList {
	results := core.ResultList{}
	for _, function := range file.Functions {
		if function.IsDefined() {
			curResults := FunctionOutOfSsaForm(function)
			results.Extend(&curResults)
		}
	}
	return results
}

// TransformFileOutOfSsaForm is the transform-pipeline adapter.
func TransformFileOutOfSsaForm(
	data *transform.TargetData,
) (*transform.TargetData, core.ResultList) {
	results := FileOutOfSsaForm(data.Code)
	return data, results
}
