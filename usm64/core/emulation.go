package usm64core

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// MARK: Error

type EmulationError interface{}

// MARK: Context

type EmulationContext struct {
	NextBlockInfo               *gen.BasicBlockInfo
	NextInstructionIndexInBlock uint
	ShouldTerminate             bool

	Registers map[string]uint64
}

func (ctx *EmulationContext) JumpToLabel(label *gen.LabelInfo) core.ResultList {
	ctx.NextBlockInfo = label.BasicBlock
	ctx.NextInstructionIndexInBlock = 0
	return core.ResultList{}
}

func (ctx *EmulationContext) ContinueToNextInstruction() core.ResultList {
	ctx.NextInstructionIndexInBlock++
	if uint(len(ctx.NextBlockInfo.Instructions)) == ctx.NextInstructionIndexInBlock {
		ctx.NextBlockInfo = ctx.NextBlockInfo.NextBlock
		ctx.NextInstructionIndexInBlock = 0
	}
	return core.ResultList{}
}

func (ctx *EmulationContext) ArgumentToValue(
	argument gen.ArgumentInfo,
) (uint64, core.ResultList) {
	switch typedArgument := argument.(type) {
	case *gen.RegisterArgumentInfo:
		name := typedArgument.Register.Name
		value, ok := ctx.Registers[name]
		if !ok {
			v := argument.Declaration()
			return 0, list.FromSingle(core.Result{{
				Type:     core.InternalErrorResult,
				Message:  "Undefined register",
				Location: &v,
			}})
		}

		return value, core.ResultList{}

	case *gen.ImmediateInfo:
		if !typedArgument.Value.IsInt64() {
			v := argument.Declaration()
			return 0, list.FromSingle(core.Result{{
				Type:     core.ErrorResult,
				Message:  "Immediate overflows 64 bits",
				Location: &v,
			}})
		}

		return uint64(typedArgument.Value.Int64()), core.ResultList{}

	case *gen.LabelArgumentInfo:
		v := argument.Declaration()
		return 0, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Expected valued argument",
			Location: &v,
		}})

	default:
		v := argument.Declaration()
		return 0, list.FromSingle(core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Unexpected argument type",
			Location: &v,
		}})
	}
}

func NewEmulationContext(
	function *gen.FunctionInfo,
) (*EmulationContext, core.ResultList) {
	return &EmulationContext{
		NextInstructionIndexInBlock: 0,
		NextBlockInfo:               function.EntryBlock,
		ShouldTerminate:             false,
		Registers:                   make(map[string]uint64),
	}, core.ResultList{}
}

type Emulateable interface {
	Emulate(ctx *EmulationContext) core.ResultList
}
