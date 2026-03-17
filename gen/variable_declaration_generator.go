package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type VariableDeclarationGenerator struct {
	ReferencedTypeGenerator FileContextGenerator[parse.TypeNode, ReferencedTypeInfo]
}

func NewVariableDeclarationGenerator() FunctionContextGenerator[parse.VariableDeclarationNode, *VariableInfo] {
	return FunctionContextGenerator[parse.VariableDeclarationNode, *VariableInfo](
		&VariableDeclarationGenerator{
			ReferencedTypeGenerator: NewReferencedTypeGenerator(),
		},
	)
}

func (g *VariableDeclarationGenerator) Generate(
	ctx *FunctionGenerationContext,
	node parse.VariableDeclarationNode,
) (*VariableInfo, core.ResultList) {
	name := NodeToSourceString(ctx.FileGenerationContext, node.Variable)

	if existing := ctx.Variables.GetVariable(name); existing != nil {
		v := node.Variable.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Variable already declared",
			Location: &v,
		}})
	}

	typ, results := g.ReferencedTypeGenerator.Generate(ctx.FileGenerationContext, node.Type)
	if !results.IsEmpty() {
		return nil, results
	}

	decl := node.Variable.View()
	variable := &VariableInfo{
		Name:        name,
		Type:        typ,
		Declaration: decl,
	}

	results = ctx.Variables.NewVariable(variable)
	return variable, results
}
