package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type FileGenerator[InstT BaseInstruction] struct {
	NamedTypeGenerator FileContextGenerator[InstT, parse.TypeDeclarationNode, *NamedTypeInfo]
	FunctionGenerator  FileContextGenerator[InstT, parse.FunctionNode, *FunctionInfo[InstT]]
}

func NewFileGenerator[InstT BaseInstruction]() FileGenerator[InstT] {
	return FileGenerator[InstT]{
		NamedTypeGenerator: NewNamedTypeGenerator[InstT](),
		FunctionGenerator:  NewFunctionGenerator[InstT](),
	}
}

func (g *FileGenerator[InstT]) createFileContext(
	ctx *GenerationContext[InstT],
	source core.SourceContext,
) *FileGenerationContext[InstT] {
	return &FileGenerationContext[InstT]{
		GenerationContext: ctx,
		SourceContext:     source,
		Types:             ctx.TypeManagerCreator(),
	}
}

func (g *FileGenerator[InstT]) generateTypesFromDeclarations(
	ctx *FileGenerationContext[InstT],
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

func (g *FileGenerator[InstT]) generateFunctions(
	ctx *FileGenerationContext[InstT],
	nodes []parse.FunctionNode,
) (functions []*FunctionInfo[InstT], results core.ResultList) {
	functions = make([]*FunctionInfo[InstT], len(nodes))
	for i, node := range nodes {
		functionInfo, curResults := g.FunctionGenerator.Generate(ctx, node)
		functions[i] = functionInfo
		results.Extend(&curResults)
	}
	return functions, results
}

func (g *FileGenerator[InstT]) Generate(
	ctx *GenerationContext[InstT],
	source core.SourceContext,
	node parse.FileNode,
) (*FileInfo[InstT], core.ResultList) {
	var results core.ResultList
	fileCtx := g.createFileContext(ctx, source)

	_, typeResults := g.generateTypesFromDeclarations(fileCtx, node.Types)
	results.Extend(&typeResults)

	functions, functionResults := g.generateFunctions(fileCtx, node.Functions)
	results.Extend(&functionResults)

	if !results.IsEmpty() {
		return nil, results
	}

	file := &FileInfo[InstT]{
		Functions: functions,
	}

	return file, core.ResultList{}
}
