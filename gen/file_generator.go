package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type FileGenerator struct {
	NamedTypeGenerator      FileContextGenerator[parse.TypeDeclarationNode, *NamedTypeInfo]
	FunctionGenerator       FileContextGenerator[parse.FunctionNode, *FunctionInfo]
	FunctionGlobalGenerator FileContextGenerator[parse.FunctionNode, GlobalInfo]
}

func NewFileGenerator() FileGenerator {
	return FileGenerator{
		NamedTypeGenerator:      NewNamedTypeGenerator(),
		FunctionGenerator:       NewFunctionGenerator(),
		FunctionGlobalGenerator: NewFunctionGlobalGenerator(),
	}
}

func (g *FileGenerator) generateTypesFromDeclarations(
	ctx *FileGenerationContext,
	nodes []parse.TypeDeclarationNode,
) (results core.ResultList) {
	for _, node := range nodes {
		_, curResults := g.NamedTypeGenerator.Generate(ctx, node)
		results.Extend(&curResults)
	}
	return results
}

func (g *FileGenerator) generateFunctionGlobals(
	ctx *FileGenerationContext,
	nodes []parse.FunctionNode,
) (results core.ResultList) {
	for _, node := range nodes {
		_, curResults := g.FunctionGlobalGenerator.Generate(ctx, node)
		results.Extend(&curResults)
	}

	return results
}

func (g *FileGenerator) generateGlobals(
	ctx *FileGenerationContext,
	node parse.FileNode,
) (results core.ResultList) {
	return g.generateFunctionGlobals(ctx, node.Functions)
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
	fileCtx := ctx.NewFileGenerationContext(source)

	results := g.generateTypesFromDeclarations(fileCtx, node.Types)
	if !results.IsEmpty() {
		return nil, results
	}

	results = g.generateGlobals(fileCtx, node)
	if !results.IsEmpty() {
		return nil, results
	}

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
