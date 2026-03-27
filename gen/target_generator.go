package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type TypedRegisterGenerator struct {
	ReferencedTypeGenerator FileContextGenerator[parse.TypeNode, ReferencedTypeInfo]
}

func NewTypedRegisterGenerator() FunctionContextGenerator[parse.TargetNode, ArgumentInfo] {
	return FunctionContextGenerator[parse.TargetNode, ArgumentInfo](
		&TypedRegisterGenerator{
			ReferencedTypeGenerator: NewReferencedTypeGenerator(),
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

// The target generator creates returns the register information that matches
// the provided target node.
//
// If the target node does not have an explicit type, and the register has not
// been defined and processed yet, the generator will return nil.
func (g *TypedRegisterGenerator) Generate(
	ctx *FunctionGenerationContext,
	node parse.TargetNode,
) (ArgumentInfo, core.ResultList) {
	registerName := NodeToSourceString(ctx.FileGenerationContext, node.Register)
	registerInfo := ctx.Registers.GetRegister(registerName)
	nodeView := node.View()

	explicitTypeProvided := node.Type != nil
	if !explicitTypeProvided {
		// If an explicit type is not provided, the best we can do is to return
		// the previously defined register information. If it is not have been
		// defined yet, we return nil here.
		if registerInfo == nil {
			return nil, core.ResultList{}
		}

		return &RegisterArgumentInfo{
			Register:    registerInfo,
			declaration: &nodeView,
		}, core.ResultList{}

	} else {
		targetType, results := g.ReferencedTypeGenerator.Generate(
			ctx.FileGenerationContext,
			*node.Type,
		)

		if !results.IsEmpty() {
			return nil, results
		}

		registerAlreadyDefined := registerInfo != nil
		if registerAlreadyDefined {
			// Register is already defined, so we just ensure that the explicit type
			// we got now matches the previously defined one.
			if !registerInfo.Type.Equal(targetType) {
				return nil, NewRegisterTypeMismatchResult(
					node.View(),
					registerInfo.Declaration,
				)
			}

			return &RegisterArgumentInfo{
				Register:    registerInfo,
				declaration: &nodeView,
			}, core.ResultList{}

		} else {
			// Register is not defined yet, so we define it now.
			registerInfo = &RegisterInfo{
				Name:        registerName,
				Type:        targetType,
				Declaration: nodeView,
			}

			results := ctx.Registers.NewRegister(registerInfo)
			if !results.IsEmpty() {
				return nil, results
			}

			return &RegisterArgumentInfo{
				Register:    registerInfo,
				declaration: &nodeView,
			}, core.ResultList{}
		}
	}
}
