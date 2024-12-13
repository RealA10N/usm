package usm64core

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Argument interface {
	String(ctx *EmulationContext) string
}

type ValuedArgument interface {
	Argument
	Value(ctx *EmulationContext) uint64
}

func NewArgument(argument gen.ArgumentInfo) (Argument, core.ResultList) {
	switch typedArgument := argument.(type) {
	case *gen.ImmediateInfo:
		return NewImmediate(*typedArgument)
	case *gen.RegisterArgumentInfo:
		return NewRegister(typedArgument.Register)
	default:
		return nil, list.FromSingle(core.Result{
			{
				Type:    core.InternalErrorResult,
				Message: "Unknown argument type",
			},
		})
	}
}
