package usm64isa

import (
	"fmt"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type FixedInstructionDefinition struct {
	Targets   core.UsmUint
	Arguments core.UsmUint
	Creator   func(
		info *gen.InstructionInfo,
	) (gen.BaseInstruction, core.ResultList)
}

func (d *FixedInstructionDefinition) InferTargetTypes(
	ctx *gen.FunctionGenerationContext,
	targets []*gen.ReferencedTypeInfo,
	arguments []*gen.ReferencedTypeInfo,
) ([]gen.ReferencedTypeInfo, core.ResultList) {
	base := ctx.Types.GetType("$64")
	if base == nil {
		return nil, list.FromSingle(core.Result{
			{
				Type:    core.InternalErrorResult,
				Message: "The $64 type is not defined",
			},
		})
	}

	inferredTargets := make([]gen.ReferencedTypeInfo, len(targets))
	for i := core.UsmUint(0); i < d.Targets; i++ {
		inferredTargets[i] = gen.ReferencedTypeInfo{Base: base}
	}

	return inferredTargets, core.ResultList{}
}

func (d *FixedInstructionDefinition) assertTargetAmount(
	targets []*gen.TargetInfo,
) core.ResultList {
	// TODO: possible overflow?
	if core.UsmUint(len(targets)) != d.Targets {
		return list.FromSingle(core.Result{
			{
				Type:    core.ErrorResult,
				Message: fmt.Sprintf("Exactly %d target(s) are allowed", d.Targets),
			},
		})
	}
	return core.ResultList{}
}

func (d *FixedInstructionDefinition) assertArgumentAmount(
	arguments []gen.ArgumentInfo,
) core.ResultList {
	// TODO: possible overflow?
	if core.UsmUint(len(arguments)) != d.Arguments {
		return list.FromSingle(core.Result{
			{
				Type:    core.ErrorResult,
				Message: fmt.Sprintf("Exactly %d argument(s) are allowed", d.Arguments),
			},
		})
	}
	return core.ResultList{}
}

func (d *FixedInstructionDefinition) assertInputLengths(
	targetInfos []*gen.TargetInfo,
	argumentInfos []gen.ArgumentInfo,
) (results core.ResultList) {
	targetResults := d.assertTargetAmount(targetInfos)
	results.Extend(&targetResults)

	argumentResults := d.assertArgumentAmount(argumentInfos)
	results.Extend(&argumentResults)

	return results
}

func (d *FixedInstructionDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := d.assertInputLengths(info.Targets, info.Arguments)
	if !results.IsEmpty() {
		return nil, results
	}

	return d.Creator(info)
}
