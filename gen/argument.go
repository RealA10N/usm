package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Info

type ArgumentInfo struct {
	// A pointer to the TypeInfo instance that corresponds to the type of the
	// register.
	Type *TypeInfo
}

type ImmediateInfo struct {
	Type  *TypeInfo
	Value core.UsmUint // TODO: more complex and complete representation of immediate structs.
}

type LabelInfo struct {
	// TODO: add location relevant information. How exactly?
	Name string
}

type GlobalInfo struct {
	Name string
	Type *TypeInfo
}

// MARK: Generator

type ArgumentGenerator[InstT BaseInstruction] struct {
	RegisterArgumentGenerator[InstT]
}

func (g *ArgumentGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.ArgumentNode,
) (*ArgumentInfo, core.ResultList) {
	switch arg := node.(type) {
	case parse.RegisterNode:
		return g.RegisterArgumentGenerator.Generate(ctx, arg)
	default:
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Unsupported argument type",
			Location: &v,
		}})
	}
}
