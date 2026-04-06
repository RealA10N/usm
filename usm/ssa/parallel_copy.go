package usmssa

import (
	"fmt"

	"alon.kr/x/usm/gen"
	usmisa "alon.kr/x/usm/usm/isa"
)

type parallelCopyEntry struct {
	dst    *gen.RegisterInfo
	srcArg gen.ArgumentInfo
}

func newMoveInstruction(dst *gen.RegisterInfo, srcArg gen.ArgumentInfo) *gen.InstructionInfo {
	instr := gen.NewEmptyInstructionInfo(nil)
	instr.SetInstruction(usmisa.NewMove())
	target := gen.NewTargetInfo(dst)
	instr.AppendTarget(&target)
	instr.AppendArgument(srcArg)
	return instr
}

// sequentializeParallelCopies converts a set of parallel copies into a correct
// sequential list of move instructions. Cycles are broken using fresh
// temporaries named "%phi_tmp_<N>".
func sequentializeParallelCopies(
	copies []parallelCopyEntry,
	function *gen.FunctionInfo,
	tmpCounter *int,
) []*gen.InstructionInfo {
	if len(copies) == 0 {
		return nil
	}

	// Build the set of destination registers.
	dstSet := make(map[*gen.RegisterInfo]bool, len(copies))
	for _, c := range copies {
		dstSet[c.dst] = true
	}

	// Pre-save any source register that is also a destination register.
	// This avoids incorrect overwrites when copies form a cycle (e.g. a swap).
	saved := make(map[*gen.RegisterInfo]*gen.RegisterInfo)
	var preInstructions []*gen.InstructionInfo

	for _, c := range copies {
		srcRegArg, ok := c.srcArg.(*gen.RegisterArgumentInfo)
		if !ok {
			continue // immediates cannot form cycles
		}
		srcReg := srcRegArg.Register
		if !dstSet[srcReg] {
			continue // this source is not overwritten by any copy
		}
		if _, alreadySaved := saved[srcReg]; alreadySaved {
			continue // already saved in a previous iteration
		}

		tmpName := fmt.Sprintf("%%phi_tmp_%d", *tmpCounter)
		*tmpCounter++
		tmpReg := gen.NewRegisterInfo(tmpName, srcReg.Type)
		function.Registers.NewRegister(tmpReg)
		saved[srcReg] = tmpReg

		preInstructions = append(preInstructions, newMoveInstruction(tmpReg, gen.NewRegisterArgumentInfo(srcReg)))
	}

	// Emit the actual copies, substituting temporaries where the source was saved.
	var copyInstructions []*gen.InstructionInfo
	for _, c := range copies {
		srcArg := c.srcArg
		if srcRegArg, ok := c.srcArg.(*gen.RegisterArgumentInfo); ok {
			if tmpReg, wasSaved := saved[srcRegArg.Register]; wasSaved {
				srcArg = gen.NewRegisterArgumentInfo(tmpReg)
			}
		}
		copyInstructions = append(copyInstructions, newMoveInstruction(c.dst, srcArg))
	}

	return append(preInstructions, copyInstructions...)
}
