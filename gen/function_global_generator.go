package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type FunctionGlobalGenerator struct{}

func NewFunctionGlobalGenerator() FileContextGenerator[parse.FunctionNode, GlobalInfo] {
	return &FunctionGlobalGenerator{}
}

func (f *FunctionGlobalGenerator) Generate(
	ctx *FileGenerationContext,
	node parse.FunctionNode,
) (GlobalInfo, core.ResultList) {
	info := &FunctionInfo{
		Name:        ViewToSourceString(ctx, node.Signature.Identifier),
		Declaration: &node.UnmanagedSourceView,
	}

	global := &FunctionGlobalInfo{
		FunctionInfo: info,
	}

	results := ctx.Globals.NewGlobal(global)
	if !results.IsEmpty() {
		return nil, results
	}

	return global, results
}
