package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
)

func ArgumentsToRegisters(
	arguments []ArgumentInfo,
) []*RegisterInfo {
	registers := []*RegisterInfo{}

	for _, arg := range arguments {
		if regArg, ok := arg.(*RegisterArgumentInfo); ok {
			registers = append(registers, regArg.Register)
		}
	}

	return registers
}

func ArgumentsToLabels(arguments []ArgumentInfo) []*LabelInfo {
	labels := []*LabelInfo{}

	for _, arg := range arguments {
		if labelArg, ok := arg.(*LabelArgumentInfo); ok {
			labels = append(labels, labelArg.Label)
		}
	}

	return labels
}

func ArgumentToType(arg ArgumentInfo) (ReferencedTypeInfo, core.ResultList) {
	switch typedArg := arg.(type) {

	case *RegisterArgumentInfo:
		return typedArg.Register.Type, core.ResultList{}

	case *ImmediateInfo:
		return typedArg.Type, core.ResultList{}

	default:
		return ReferencedTypeInfo{}, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected a typed argument",
				Location: arg.Declaration(),
			},
		})
	}
}

func ArgumentsToTypes(
	arguments []ArgumentInfo,
) ([]ReferencedTypeInfo, core.ResultList) {
	results := core.ResultList{}
	types := make([]ReferencedTypeInfo, len(arguments))

	for i, arg := range arguments {
		typ, curResults := ArgumentToType(arg)
		results.Extend(&curResults)
		types[i] = typ
	}

	if !results.IsEmpty() {
		return nil, results
	}

	return types, core.ResultList{}
}

func ArgumentToLabel(arg ArgumentInfo) (*LabelInfo, core.ResultList) {
	if labelArg, ok := arg.(*LabelArgumentInfo); ok {
		return labelArg.Label, core.ResultList{}
	}

	return nil, list.FromSingle(core.Result{
		{
			Type:     core.ErrorResult,
			Message:  "Expected a label argument",
			Location: arg.Declaration(),
		},
	})
}

func TargetsToRegisters(targets []*TargetInfo) []*RegisterInfo {
	registers := []*RegisterInfo{}

	for _, target := range targets {
		registers = append(registers, target.Register)
	}

	return registers
}

func TargetToType(target *TargetInfo) ReferencedTypeInfo {
	return target.Register.Type
}
