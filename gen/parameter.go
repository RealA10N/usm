package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type ParameterGenerator[InstT BaseInstruction] struct {
	ReferencedTypeGenerator Generator[InstT, parse.TypeNode, ReferencedTypeInfo]
}

func NewParameterGenerator[InstT BaseInstruction]() Generator[InstT, parse.ParameterNode, *RegisterInfo] {
	return Generator[InstT, parse.ParameterNode, *RegisterInfo](
		&ParameterGenerator[InstT]{
			ReferencedTypeGenerator: NewReferencedTypeGenerator[InstT](),
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
	ctx *GenerationContext[InstT],
	node parse.ParameterNode,
) (*RegisterInfo, core.ResultList) {
	results := core.ResultList{}

	typeInfo, typeResults := g.ReferencedTypeGenerator.Generate(ctx, node.Type)
	results.Extend(&typeResults)

	registerName := nodeToSourceString(ctx, node.Register)
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
