package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type ParameterGenerator[InstT BaseInstruction] struct {
	ReferencedTypeGenerator FileContextGenerator[parse.TypeNode, ReferencedTypeInfo]
}

func NewParameterGenerator[InstT BaseInstruction]() FunctionContextGenerator[InstT, parse.ParameterNode, *RegisterInfo] {
	return FunctionContextGenerator[InstT, parse.ParameterNode, *RegisterInfo](
		&ParameterGenerator[InstT]{
			ReferencedTypeGenerator: NewReferencedTypeGenerator(),
		},
	)
}

func NewRegisterAlreadyDefinedResult(
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
func (g *ParameterGenerator[InstT]) Generate(
	ctx *FunctionGenerationContext[InstT],
	node parse.ParameterNode,
) (*RegisterInfo, core.ResultList) {
	results := core.ResultList{}

	typeInfo, typeResults := g.ReferencedTypeGenerator.Generate(
		ctx.FileGenerationContext,
		node.Type,
	)
	results.Extend(&typeResults)

	registerName := nodeToSourceString(ctx.FileGenerationContext, node.Register)
	registerInfo := ctx.Registers.GetRegister(registerName)
	if registerInfo != nil {
		registerResults := NewRegisterAlreadyDefinedResult(
			node.View(),
			registerInfo.Declaration,
		)
		results.Extend(&registerResults)
	}

	if !results.IsEmpty() {
		return nil, results
	}

	registerInfo = &RegisterInfo{
		Name:        registerName,
		Type:        typeInfo,
		Declaration: node.View(),
	}

	result := ctx.Registers.NewRegister(registerInfo)
	if result != nil {
		return nil, list.FromSingle(result)
	}

	return registerInfo, core.ResultList{}
}
