package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Context

// This structure is the most broad level of generation context.
// It contains information that is used in different parts of the compilation,
// and which are essential across the whole pipeline.
//
// It mainly contains information about the compilation target.
type GenerationContext struct {
	// The size of a pointer type in the current target architecture.
	//
	// TODO: I'm not sure that we need this information in this step of the
	//   compilation. Can we compile without it, and leave it to the ISA?
	PointerSize core.UsmUint
}

type FileGenerationContext struct {
	*GenerationContext

	// The source code of the file that we are currently processing.
	core.SourceContext

	// A type manager that contains all types defined in the file,
	// and that the compiler can use to create new types when in processes
	// new type definitions.
	Types TypeManager

	// TODO: add globals, variables, constants.
}

type LabelGenerationContext struct {
	*FileGenerationContext

	// The index of the instruction which is currently being iterated upon.
	//
	// Used in the pass before we generate the instruction instances, to
	// go over the labels in a function and give each label a corresponding
	// instruction index.
	CurrentInstructionIndex core.UsmUint
}

type FunctionGenerationContext[InstT BaseInstruction] struct {
	*FileGenerationContext

	// An instruction (definition) manager, which contains all instruction
	// definitions that are supported in the current architecture (ISA).
	//
	// When processing a new instruction from the source code, the compiler
	// talks with the instruction manager, retrieves the relevant instruction
	// definition, and uses it to farther process the instruction.
	Instructions InstructionManager[InstT]

	// A register manager, which contains all active registers in the function.
	// The compiler can query register information from the register manager
	// when it encounters registers while processing, and can create new register
	// information structures and pass them to the manager, which stores them.
	Registers RegisterManager

	Labels LabelManager
}

// MARK: Generator

type Generator[NodeT parse.Node, InfoT any] interface {
	Generate(
		ctx *GenerationContext,
		node NodeT,
	) (info InfoT, results core.ResultList)
}

type FileContextGenerator[NodeT parse.Node, InfoT any] interface {
	Generate(
		ctx *FileGenerationContext,
		node NodeT,
	) (info InfoT, results core.ResultList)
}

type LabelContextGenerator[NodeT parse.Node, InfoT any] interface {
	Generate(
		ctx *LabelGenerationContext,
		node NodeT,
	) (info InfoT, results core.ResultList)
}

type FunctionContextGenerator[InstT BaseInstruction, NodeT parse.Node, InfoT any] interface {
	Generate(
		ctx *FunctionGenerationContext[InstT],
		node NodeT,
	) (info InfoT, results core.ResultList)
}

// MARK: Utils

func viewToSourceString(
	ctx *FileGenerationContext,
	view core.UnmanagedSourceView,
) string {
	return string(view.Raw(ctx.SourceContext))
}

func nodeToSourceString(
	ctx *FileGenerationContext,
	node parse.Node,
) string {
	return viewToSourceString(ctx, node.View())
}
