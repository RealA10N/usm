package usm64core

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Argument interface {
	String(ctx *EmulationContext) string
	Declaration() core.UnmanagedSourceView
}

type ValuedArgument interface {
	Argument
	Value(ctx *EmulationContext) uint64
}

func ArgumentToValuedArgument(arg Argument) (ValuedArgument, core.ResultList) {
	valued, ok := arg.(ValuedArgument)
	if !ok {
		v := arg.Declaration()
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected valued argument",
				Location: &v,
			},
		})
	}

	return valued, core.ResultList{}
}

func ArgumentToLabel(argument Argument) (Label, core.ResultList) {
	label, ok := argument.(Label)
	if !ok {
		v := argument.Declaration()
		return Label{}, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Expected label argument",
				Location: &v,
			},
		})
	}

	return label, core.ResultList{}
}

func NewArgument(argument gen.ArgumentInfo) (Argument, core.ResultList) {
	switch typedArgument := argument.(type) {
	case *gen.ImmediateInfo:
		return NewImmediate(*typedArgument)
	case *gen.RegisterArgumentInfo:
		return NewRegister(typedArgument)
	case *gen.LabelArgumentInfo:
		return NewLabel(*typedArgument)
	default:
		return nil, list.FromSingle(core.Result{
			{
				Type:    core.InternalErrorResult,
				Message: "Unknown argument type",
			},
		})
	}
}
