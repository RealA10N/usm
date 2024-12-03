package gen

import (
	"math/big"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

type ImmediateInfo struct {
	Type  *TypeInfo
	Value *big.Int // TODO: Add floating types
	// TODO: more complex and complete representation of immediate structs.
}

func (i *ImmediateInfo) GetType() *TypeInfo {
	return i.Type
}

// MARK: Generator

type ImmediateArgumentGenerator[InstT BaseInstruction] struct{}

func (g *ImmediateArgumentGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.ImmediateNode,
) (ArgumentInfo, core.ResultList) {
	typeName := string(node.Type.Identifier.Raw(ctx.SourceContext))
	typ := ctx.Types.GetType(typeName)
	if typ == nil {
		return nil, list.FromSingle(NewUndefinedTypeResult(node.View()))
	}

	immediate, ok := node.Value.(parse.ImmediateFinalValueNode)
	if !ok || len(node.Type.Decorators) > 0 {
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Complex immediate values are not supported yet",
			Location: &v,
		}})
	}

	immediate.Start += 1 // to skip the '#' character
	valueStr := string(immediate.Raw(ctx.SourceContext))
	value, ok := new(big.Int).SetString(valueStr, 0)
	if !ok {
		v := immediate.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Invalid immediate value",
			Location: &v,
		}})
	}

	info := ImmediateInfo{
		Type:  typ,
		Value: value,
	}

	return &info, core.ResultList{}
}
