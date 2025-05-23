package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type ParameterGenerator struct {
	ReferencedTypeGenerator FileContextGenerator[parse.TypeNode, ReferencedTypeInfo]
}

func NewParameterGenerator() FunctionContextGenerator[parse.ParameterNode, *RegisterInfo] {
	return FunctionContextGenerator[parse.ParameterNode, *RegisterInfo](
		&ParameterGenerator{
			ReferencedTypeGenerator: NewReferencedTypeGenerator(),
		},
	)
}

func newRegisterAlreadyDefinedResult(
	NewDeclaration core.UnmanagedSourceView,
	FirstDeclaration core.UnmanagedSourceView,
) core.ResultList {
	return list.FromSingle(core.Result{
		{
			Type:     core.ErrorResult,
			Message:  "Register already defined",
			Location: &NewDeclaration,
		},
		{
			Type:     core.HintResult,
			Message:  "Previous definition here",
			Location: &FirstDeclaration,
		},
	})
}

// Asserts that a register with the same name does not exist yet,
// creates the new register, registers it to the register manager,
// and returns the unique register info structure pointer.
func (g *ParameterGenerator) Generate(
	ctx *FunctionGenerationContext,
	node parse.ParameterNode,
) (*RegisterInfo, core.ResultList) {
	results := core.ResultList{}

	typeInfo, typeResults := g.ReferencedTypeGenerator.Generate(
		ctx.FileGenerationContext,
		node.Type,
	)

	results.Extend(&typeResults)
	if !results.IsEmpty() {
		return nil, results
	}

	registerName := NodeToSourceString(ctx.FileGenerationContext, node.Register)
	registerInfo := ctx.Registers.GetRegister(registerName)
	if registerInfo != nil {
		return registerInfo, core.ResultList{}
	}

	registerInfo = &RegisterInfo{
		Name:        registerName,
		Type:        typeInfo,
		Declaration: node.View(),
	}

	results = ctx.Registers.NewRegister(registerInfo)
	if !results.IsEmpty() {
		return nil, results
	}

	return registerInfo, core.ResultList{}
}
