package usm64core

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

func instructionInfoToInstruction(info *gen.InstructionInfo) (Instruction, core.ResultList) {
	if inst, ok := info.Instruction.(Instruction); !ok {
		return nil, list.FromSingle(core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Invalid instruction type",
			Location: info.Declaration,
		}})
	} else {
		return inst, core.ResultList{}
	}
}

type Emulator struct{}

func (Emulator) Emulate(
	function *gen.FunctionInfo,
) core.ResultList {
	ctx, results := NewEmulationContext(function)
	if !results.IsEmpty() {
		return results
	}

	for !ctx.ShouldTerminate {
		instInfo := ctx.NextBlockInfo.Instructions[ctx.NextInstructionIndexInBlock]
		instruction, results := instructionInfoToInstruction(instInfo)
		if !results.IsEmpty() {
			return results
		}

		results = instruction.Emulate(ctx)
		if !results.IsEmpty() {
			return results
		}
	}

	return core.ResultList{}
}

func NewEmulator() Emulator {
	return Emulator{}
}
