package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Generator

type TargetGenerator[InstT BaseInstruction] struct {
	ReferencedTypeGenerator FileContextGenerator[InstT, parse.TypeNode, ReferencedTypeInfo]
}

func NewTargetGenerator[InstT BaseInstruction]() FunctionContextGenerator[InstT, parse.TargetNode, partialRegisterInfo] {
	return FunctionContextGenerator[InstT, parse.TargetNode, partialRegisterInfo](
		&TargetGenerator[InstT]{
			ReferencedTypeGenerator: NewReferencedTypeGenerator[InstT](),
		},
	)
}

func NewRegisterTypeMismatchResult(
	NewDeclaration core.UnmanagedSourceView,
	FirstDeclaration core.UnmanagedSourceView,
) core.ResultList {
	return list.FromSingle(core.Result{
		{
			Type:     core.ErrorResult,
			Message:  "Explicit register type does not match previous declaration",
			Location: &NewDeclaration,
		},
		{
			Type:     core.HintResult,
			Message:  "Previous declaration here",
			Location: &FirstDeclaration,
		},
	})
}

func (g *TargetGenerator[InstT]) Generate(
	ctx *FunctionGenerationContext[InstT],
	node parse.TargetNode,
) (partialRegisterInfo, core.ResultList) {
	var explicitType *ReferencedTypeInfo

	// if an explicit type is provided to the target, get the type info.
	explicitTypeProvided := node.Type != nil
	if explicitTypeProvided {
		explicitTypeValue, results := g.ReferencedTypeGenerator.Generate(
			ctx.FileGenerationContext,
			*node.Type,
		)
		if !results.IsEmpty() {
			return partialRegisterInfo{}, results
		}

		explicitType = &explicitTypeValue
	}

	registerName := nodeToSourceString(ctx.FileGenerationContext, node.Register)
	registerInfo := ctx.Registers.GetRegister(registerName)

	registerAlreadyDefined := registerInfo != nil
	if registerAlreadyDefined {
		if explicitType != nil {
			// ensure explicit type matches the previously declared one.
			if !explicitType.Equal(registerInfo.Type) {
				return partialRegisterInfo{}, NewRegisterTypeMismatchResult(
					node.View(),
					registerInfo.Declaration,
				)
			}
		}

		// all checks passed; return previously defined register type.
		return registerInfo.toPartialRegisterInfo(), core.ResultList{}

	} else {
		// this is the first appearance of the register; if the type is provided
		// explicitly, use it. otherwise, there is no way to know the type of
		// the target register at this.
		// the type and register will be finalized when the instruction is built,
		// and only then it is added to the register manager.
		return partialRegisterInfo{
			Name:        registerName,
			Type:        explicitType,
			Declaration: node.View(),
		}, core.ResultList{}
	}
}
