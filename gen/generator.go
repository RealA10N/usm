package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type ArchInfo struct {
	PointerSize core.UsmUint // The size of a pointer in bytes.
}

// MARK: Context

// A context object that is initialized empty, but gets propagated and filled
// with information as the code generation process continues, while iterating
// over the AST nodes.
type GenerationContext[InstT BaseInstruction] struct {
	ArchInfo
	core.SourceContext

	Instructions InstructionManager[InstT]
	Types        TypeManager
	Registers    RegisterManager
	// TODO: add globals, functions, etc.
}

// MARK: Generator

type Generator[InstT BaseInstruction, NodeT parse.Node, InfoT any] interface {
	Generate(
		ctx *GenerationContext[InstT],
		node NodeT,
	) (info InfoT, results core.ResultList)
}

func viewToSourceString[InstT BaseInstruction](
	ctx *GenerationContext[InstT],
	view core.UnmanagedSourceView,
) string {
	return string(view.Raw(ctx.SourceContext))
}

func nodeToSourceString[InstT BaseInstruction](
	ctx *GenerationContext[InstT],
	node parse.Node,
) string {
	return viewToSourceString(ctx, node.View())
}
