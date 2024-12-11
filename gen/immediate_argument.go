package gen

import (
	"math/big"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

type ImmediateInfo struct {
	Type  ReferencedTypeInfo
	Value *big.Int // TODO: Add floating types
	// TODO: more complex and complete representation of immediate structs.
}

func (i *ImmediateInfo) GetType() ReferencedTypeInfo {
	return i.Type
}

// MARK: Generator

type ImmediateArgumentGenerator[InstT BaseInstruction] struct {
	ReferencedTypeGenerator Generator[InstT, parse.TypeNode, ReferencedTypeInfo]
}

func NewImmediateArgumentGenerator[InstT BaseInstruction]() Generator[InstT, parse.ImmediateNode, ArgumentInfo] {
	return Generator[InstT, parse.ImmediateNode, ArgumentInfo](
		&ImmediateArgumentGenerator[InstT]{
			ReferencedTypeGenerator: NewReferencedTypeGenerator[InstT](),
		},
	)
}

func (g *ImmediateArgumentGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.ImmediateNode,
) (ArgumentInfo, core.ResultList) {
	typeInfo, results := g.ReferencedTypeGenerator.Generate(ctx, node.Type)
	if !results.IsEmpty() {
		return nil, results
	}

	immediate, ok := node.Value.(parse.ImmediateFinalValueNode)
	if !ok {
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Complex immediate values are not supported yet",
			Location: &v,
		}})
	}

	immediate.Start += 1 // to skip the '#' character
	valueStr := nodeToSourceString(ctx, immediate)
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
		Type:  typeInfo,
		Value: value,
	}

	return &info, core.ResultList{}
}
