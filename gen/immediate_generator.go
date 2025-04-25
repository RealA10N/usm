package gen

import (
	"math/big"

	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type ImmediateArgumentGenerator struct {
	ReferencedTypeGenerator FileContextGenerator[parse.TypeNode, ReferencedTypeInfo]
}

func NewImmediateArgumentGenerator() InstructionContextGenerator[parse.ImmediateNode, ArgumentInfo] {
	return InstructionContextGenerator[parse.ImmediateNode, ArgumentInfo](
		&ImmediateArgumentGenerator{
			ReferencedTypeGenerator: NewReferencedTypeGenerator(),
		},
	)
}

func (g *ImmediateArgumentGenerator) Generate(
	ctx *InstructionGenerationContext,
	node parse.ImmediateNode,
) (ArgumentInfo, core.ResultList) {
	typeInfo, results := g.ReferencedTypeGenerator.Generate(
		ctx.FileGenerationContext,
		node.Type,
	)
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
	valueStr := NodeToSourceString(ctx.FileGenerationContext, immediate)
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
		Type:        typeInfo,
		Value:       value,
		declaration: node.View(),
	}

	return &info, core.ResultList{}
}
