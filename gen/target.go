package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Generator

type TargetGenerator[InstT BaseInstruction] struct{}

func (g *TargetGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.TargetNode,
) (*TypeInfo, core.ResultList) {

	// if an explicit type is provided to the target, get the type info.
	var explicitType *TypeInfo
	if node.Type != nil {
		explicitTypeName := string(node.Type.Identifier.Raw(ctx.SourceContext))
		explicitType = ctx.Types.GetType(explicitTypeName)
	}

	registerName := string(node.Register.Raw(ctx.SourceContext))
	registerInfo := ctx.Registers.GetRegister(registerName)

	if registerInfo != nil {
		// register is already previously defined
		if explicitType != nil {
			// ensure explicit type matches the previously declared one.
			if explicitType != registerInfo.Type {
				return nil, list.FromSingle(
					NewRegisterTypeMismatchResult(
						node.View(),
						registerInfo.Declaration,
					),
				)
			}
		}

		// all checks passed; return previously defined register type.
		return registerInfo.Type, core.ResultList{}

	} else {
		// this is the first appearance of the register; if the type is provided
		// explicitly, use it. otherwise, there is no way to know the type of
		// the target register at this.
		// the type and register will be finalized when the instruction is built,
		// and only then it is added to the register manager.
		return explicitType, core.ResultList{}
	}
}
