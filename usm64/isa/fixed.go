package usm64isa

import (
	"fmt"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	usm64core "alon.kr/x/usm/usm64/core"
)

type FixedInstructionDefinition struct {
	Targets   core.UsmUint
	Arguments core.UsmUint
	Creator   func(
		targets []usm64core.Register,
		arguments []usm64core.Argument,
	) (usm64core.Instruction, core.ResultList)
}

func (d *FixedInstructionDefinition) InferTargetTypes(
	ctx *gen.FunctionGenerationContext[usm64core.Instruction],
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
	targets []*gen.RegisterArgumentInfo,
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
	targetInfos []*gen.RegisterArgumentInfo,
	argumentInfos []gen.ArgumentInfo,
) (results core.ResultList) {
	targetResults := d.assertTargetAmount(targetInfos)
	results.Extend(&targetResults)

	argumentResults := d.assertArgumentAmount(argumentInfos)
	results.Extend(&argumentResults)

	return results
}

func (d *FixedInstructionDefinition) createRegisters(
	registerInfos []*gen.RegisterArgumentInfo,
) (registers []usm64core.Register, results core.ResultList) {
	registers = make([]usm64core.Register, len(registerInfos))
	for i, register := range registerInfos {
		target, curResults := usm64core.NewRegister(register)
		results.Extend(&curResults)
		registers[i] = target
	}
	return
}

func (d *FixedInstructionDefinition) createArguments(
	argumentInfos []gen.ArgumentInfo,
) (arguments []usm64core.Argument, results core.ResultList) {
	arguments = make([]usm64core.Argument, len(argumentInfos))
	for i, argument := range argumentInfos {
		argument, curResults := usm64core.NewArgument(argument)
		results.Extend(&curResults)
		arguments[i] = argument
	}
	return
}

func (d *FixedInstructionDefinition) BuildInstruction(
	targetInfos []*gen.RegisterArgumentInfo,
	argumentInfos []gen.ArgumentInfo,
) (usm64core.Instruction, core.ResultList) {
	results := d.assertInputLengths(targetInfos, argumentInfos)
	if !results.IsEmpty() {
		return nil, results
	}

	targets, results := d.createRegisters(targetInfos)
	if !results.IsEmpty() {
		return nil, results
	}

	arguments, results := d.createArguments(argumentInfos)
	if !results.IsEmpty() {
		return nil, results
	}

	return d.Creator(targets, arguments)
}
