package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	"alon.kr/x/macho/load/section64"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Bl struct {
	Target *gen.FunctionInfo
}

func (b Bl) Operator() string {
	return "bl"
}

func (b Bl) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	// TODO: add an analysis to check if the target function is a no-return
	// function.
	return gen.StepInfo{PossibleContinue: true}, core.ResultList{}
}

func (b Bl) registerRelocation(ctx *aarch64codegen.InstructionCodegenContext) {
	relocation := section64.RelocationBuilder{
		Address:                uint32(ctx.InstructionOffsetInFile()),
		SymbolIndex:            ctx.FunctionIndices[b.Target],
		IsRelocationPcRelative: true,
		Length:                 section64.RelocationLengthLong,
		IsRelocationExtern:     true,
		Type:                   section64.RelocationTypeArm64Branch26,
	}

	ctx.Relocations = append(ctx.Relocations, relocation)
}

func (b Bl) Generate(
	ctx *aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	targetOffset, ok := ctx.FunctionOffsets[b.Target]
	if !ok {
		// Target function is not defined: we add a relocation to the symbol
		// and let the linker resolve it.
		b.registerRelocation(ctx)
		return instructions.BL(0), core.ResultList{}
	}

	currentOffset := ctx.InstructionOffsetInFile()
	offset, err := aarch64translation.Uint64DiffToOffset26Align4(targetOffset, currentOffset)

	if err != nil {
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Invalid branch offset (arbitrary large offsets are not yet supported)",
				Location: ctx.Declaration,
			},
			{
				Type:    core.DebugResult,
				Message: err.Error(),
			},
		})
	}

	return instructions.BL(offset), core.ResultList{}
}

type BlDefinition struct{}

func (BlDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := core.ResultList{}

	curResults := aarch64translation.AssertArgumentsExactly(info, 1)
	results.Extend(&curResults)

	curResults = aarch64translation.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	target, curResults := aarch64translation.ArgumentToFunctionInfo(info.Arguments[0])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	return Bl{Target: target}, core.ResultList{}
}

func NewBlInstructionDefinition() gen.InstructionDefinition {
	return BlDefinition{}
}
