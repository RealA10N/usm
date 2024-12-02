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

type ArgumentGenerator[InstT BaseInstruction] struct{}

func (g *ArgumentGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.ArgumentNode,
) (*ArgumentInfo, core.ResultList) {
	switch arg := node.(type) {

	case parse.RegisterNode:
		// TODO: duplicated code: make function that extracts register name from node.
		registerName := string(arg.Raw(ctx.SourceContext))
		registerInfo := ctx.Registers.GetRegister(registerName)

		if registerInfo == nil {
			v := node.View()
			return nil, list.FromSingle(core.Result{{
				Type:     core.ErrorResult,
				Message:  "Undefined register used as argument",
				Location: &v,
			}})
		}

		argumentInfo := ArgumentInfo{
			Type: registerInfo.Type,
		}

		return &argumentInfo, core.ResultList{}

	default:
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Unsupported argument type",
			Location: &v,
		}})
	}
}
