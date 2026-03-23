package usmssa

import (
	"fmt"
	"slices"

	"alon.kr/x/usm/gen"
	ssaopt "alon.kr/x/usm/opt/ssa"
	usmisa "alon.kr/x/usm/usm/isa"
)

func newJumpInstruction(target *gen.LabelInfo) *gen.InstructionInfo {
	instr := gen.NewEmptyInstructionInfo(nil)
	instr.SetInstruction(usmisa.NewJump())
	instr.AppendArgument(gen.NewLabelArgumentInfo(target))
	return instr
}

// findBlockBefore returns the block whose NextBlock is target, or nil.
func findBlockBefore(function *gen.FunctionInfo, target *gen.BasicBlockInfo) *gen.BasicBlockInfo {
	for block := function.EntryBlock; block != nil; block = block.NextBlock {
		if block.NextBlock == target {
			return block
		}
	}
	return nil
}

// splitCriticalEdge splits the critical edge pred→succ by inserting a new
// intermediate block. The new block is named ".ssa_split_<counter>".
func splitCriticalEdge(
	function *gen.FunctionInfo,
	pred *gen.BasicBlockInfo,
	succ *gen.BasicBlockInfo,
	splitCounter *int,
) {
	// Create split block with a deterministic label.
	splitLabelName := fmt.Sprintf(".ssa_split_%d", *splitCounter)
	*splitCounter++
	splitLabel := &gen.LabelInfo{Name: splitLabelName}
	function.Labels.NewLabel(splitLabel)

	splitBlock := gen.NewEmptyBasicBlockInfo(function)
	splitBlock.SetLabel(splitLabel)
	splitBlock.AppendInstruction(newJumpInstruction(succ.Label))

	// Determine whether the edge is represented as an explicit label argument
	// in the terminator, or as an implicit fall-through (pred.NextBlock == succ).
	terminator := pred.Instructions[len(pred.Instructions)-1]
	isExplicitBranch := false
	for i, arg := range terminator.Arguments {
		if labelArg, ok := arg.(*gen.LabelArgumentInfo); ok {
			if labelArg.Label.BasicBlock == succ {
				terminator.Arguments[i] = gen.NewLabelArgumentInfo(splitLabel)
				isExplicitBranch = true
				break
			}
		}
	}

	// Insert split block in the linked list.
	if isExplicitBranch {
		// Insert split block right before succ so the linear order is sensible.
		blockBeforeSucc := findBlockBefore(function, succ)
		if blockBeforeSucc != nil {
			blockBeforeSucc.AppendBasicBlock(splitBlock)
		} else {
			// succ is the entry block; insert after pred as a fallback.
			pred.AppendBasicBlock(splitBlock)
		}
	} else {
		// Fall-through edge: insert split block between pred and succ so that
		// pred's fall-through now goes to split block → succ.
		pred.AppendBasicBlock(splitBlock)
	}

	// Update phi instructions in succ: replace references to pred with split block.
	for _, instr := range succ.Instructions {
		if _, isPhi := instr.Definition.(ssaopt.PhiInstructionDefinition); !isPhi {
			break
		}
		for i, arg := range instr.Arguments {
			if labelArg, ok := arg.(*gen.LabelArgumentInfo); ok {
				if labelArg.Label.BasicBlock == pred {
					instr.Arguments[i] = gen.NewLabelArgumentInfo(splitLabel)
				}
			}
		}
	}

	// Update CFG edges: pred→succ becomes pred→split, split→succ.
	pred.ForwardEdges = slices.DeleteFunc(
		pred.ForwardEdges,
		func(b *gen.BasicBlockInfo) bool { return b == succ },
	)
	succ.BackwardEdges = slices.DeleteFunc(
		succ.BackwardEdges,
		func(b *gen.BasicBlockInfo) bool { return b == pred },
	)
	pred.AppendForwardEdge(splitBlock)
	splitBlock.AppendForwardEdge(succ)
}

// splitCriticalEdgesInFunction splits all critical edges in the function.
// A critical edge is an edge from a block with multiple successors to a block
// with multiple predecessors.
func splitCriticalEdgesInFunction(function *gen.FunctionInfo, splitCounter *int) {
	blocks := function.CollectBasicBlocks()
	for _, pred := range blocks {
		if len(pred.ForwardEdges) <= 1 {
			continue
		}
		// Snapshot forward edges before modifying.
		succs := slices.Clone(pred.ForwardEdges)
		for _, succ := range succs {
			if len(succ.BackwardEdges) <= 1 {
				continue
			}
			splitCriticalEdge(function, pred, succ, splitCounter)
		}
	}
}
