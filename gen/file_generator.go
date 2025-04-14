package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type FileGenerator struct {
	NamedTypeGenerator FileContextGenerator[parse.TypeDeclarationNode, *NamedTypeInfo]
	FunctionGenerator  FileContextGenerator[parse.FunctionNode, *FunctionInfo]
}

func NewFileGenerator() FileGenerator {
	return FileGenerator{
		NamedTypeGenerator: NewNamedTypeGenerator(),
		FunctionGenerator:  NewFunctionGenerator(),
	}
}

func CreateFileContext(
	ctx *GenerationContext,
	source core.SourceContext,
) *FileGenerationContext {
	return &FileGenerationContext{
		GenerationContext: ctx,
		SourceContext:     source,
		Types:             ctx.TypeManagerCreator(),
	}
}

func (g *FileGenerator) generateTypesFromDeclarations(
	ctx *FileGenerationContext,
	nodes []parse.TypeDeclarationNode,
) (types []*NamedTypeInfo, results core.ResultList) {
	types = make([]*NamedTypeInfo, len(nodes))
	for i, node := range nodes {
		typeInfo, curResults := g.NamedTypeGenerator.Generate(ctx, node)
		types[i] = typeInfo
		results.Extend(&curResults)
	}
	return types, results
}

func (g *FileGenerator) generateFunctions(
	ctx *FileGenerationContext,
	nodes []parse.FunctionNode,
) (functions []*FunctionInfo, results core.ResultList) {
	functions = make([]*FunctionInfo, len(nodes))
	for i, node := range nodes {
		functionInfo, curResults := g.FunctionGenerator.Generate(ctx, node)
		functions[i] = functionInfo
		results.Extend(&curResults)
	}
	return functions, results
}

func (g *FileGenerator) Generate(
	ctx *GenerationContext,
	source core.SourceContext,
	node parse.FileNode,
) (*FileInfo, core.ResultList) {
	var results core.ResultList
	fileCtx := CreateFileContext(ctx, source)

	_, typeResults := g.generateTypesFromDeclarations(fileCtx, node.Types)
	results.Extend(&typeResults)

	functions, functionResults := g.generateFunctions(fileCtx, node.Functions)
	results.Extend(&functionResults)

	if !results.IsEmpty() {
		return nil, results
	}

	file := &FileInfo{
		Functions: functions,
	}

	return file, core.ResultList{}
}
