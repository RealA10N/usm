package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Generator

type TargetGenerator[InstT BaseInstruction] struct {
	ReferencedTypeGenerator Generator[InstT, parse.TypeNode, *ReferencedTypeInfo]
}

func NewTargetGenerator[InstT BaseInstruction]() Generator[InstT, parse.TargetNode, *ReferencedTypeInfo] {
	return Generator[InstT, parse.TargetNode, *ReferencedTypeInfo](
		&TargetGenerator[InstT]{
			ReferencedTypeGenerator: NewReferencedTypeGenerator[InstT](),
		},
	)
}

func (g *TargetGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	node parse.TargetNode,
) (*ReferencedTypeInfo, core.ResultList) {
	var explicitType *ReferencedTypeInfo

	// if an explicit type is provided to the target, get the type info.
	if node.Type != nil {
		var results core.ResultList
		explicitType, results = g.ReferencedTypeGenerator.Generate(ctx, *node.Type)
		if !results.IsEmpty() {
			return nil, results
		}
	}

	registerName := string(node.Register.Raw(ctx.SourceContext))
	registerInfo := ctx.Registers.GetRegister(registerName)

	if registerInfo != nil {
		// register is already previously defined
		if explicitType != nil {
			// ensure explicit type matches the previously declared one.
			if !explicitType.Equals(registerInfo.Type) {
				return nil, list.FromSingle(
					NewRegisterTypeMismatchResult(
						node.View(),
						registerInfo.Declaration,
					),
				)
			}
		}

		// all checks passed; return previously defined register type.
		return &registerInfo.Type, core.ResultList{}

	} else {
		// this is the first appearance of the register; if the type is provided
		// explicitly, use it. otherwise, there is no way to know the type of
		// the target register at this.
		// the type and register will be finalized when the instruction is built,
		// and only then it is added to the register manager.
		return explicitType, core.ResultList{}
	}
}
