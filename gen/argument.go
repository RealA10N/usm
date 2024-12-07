package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

type ArgumentInfo interface {
	// A pointer to the TypeInfo instance that corresponds to the type of the
	// register.
	GetType() *NamedTypeInfo
}

// MARK: Generator

type ArgumentGenerator[InstT BaseInstruction] struct {
	RegisterArgumentGenerator[InstT]
	ImmediateArgumentGenerator[InstT]
}

func (g *ArgumentGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.ArgumentNode,
) (ArgumentInfo, core.ResultList) {
	switch typedNode := node.(type) {
	case parse.RegisterNode:
		return g.RegisterArgumentGenerator.Generate(ctx, typedNode)
	case parse.ImmediateNode:
		return g.ImmediateArgumentGenerator.Generate(ctx, typedNode)
	default:
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Unsupported argument type",
			Location: &v,
		}})
	}
}
