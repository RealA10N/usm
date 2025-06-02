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
	gen.NonBranchingInstruction
}

func NewBl() gen.InstructionDefinition {
	return Bl{}
}

func (b Bl) Operator(*gen.InstructionInfo) string {
	return "bl"
}

func (b Bl) Target(
	info *gen.InstructionInfo,
) (*gen.FunctionInfo, core.ResultList) {
	results := gen.AssertArgumentsExactly(info, 1)
	if !results.IsEmpty() {
		return nil, results
	}

	target, results := aarch64translation.ArgumentToFunctionInfo(info.Arguments[0])
	if !results.IsEmpty() {
		return nil, results
	}

	return target, core.ResultList{}
}

func (b Bl) Validate(info *gen.InstructionInfo) core.ResultList {
	results := core.ResultList{}

	_, curResults := b.Target(info)
	results.Extend(&curResults)

	curResults = gen.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}

func (b Bl) registerRelocation(
	ctx *aarch64codegen.InstructionCodegenContext,
) core.ResultList {
	target, results := b.Target(ctx.InstructionInfo)
	if !results.IsEmpty() {
		return results
	}

	relocation := section64.RelocationBuilder{
		Address:                uint32(ctx.InstructionOffsetInFile()),
		SymbolIndex:            ctx.FunctionIndices[target],
		IsRelocationPcRelative: true,
		Length:                 section64.RelocationLengthLong,
		IsRelocationExtern:     true,
		Type:                   section64.RelocationTypeArm64Branch26,
	}

	ctx.Relocations = append(ctx.Relocations, relocation)
	return core.ResultList{}
}

func (b Bl) Codegen(
	ctx *aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	target, results := b.Target(ctx.InstructionInfo)
	if !results.IsEmpty() {
		return nil, results
	}

	targetOffset, ok := ctx.FunctionOffsets[target]
	if !ok {
		// Target function is not defined: we add a relocation to the symbol
		// and let the linker resolve it.
		results = b.registerRelocation(ctx)
		if !results.IsEmpty() {
			return nil, results
		}

		return instructions.NewBl(0), core.ResultList{}
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

	return instructions.NewBl(offset), core.ResultList{}
}
