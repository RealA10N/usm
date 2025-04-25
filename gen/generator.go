package gen

import (
	"math/big"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

// MARK: Context

type ManagerCreators struct {
	RegisterManagerCreator func(*FileGenerationContext) RegisterManager
	LabelManagerCreator    func(*FileGenerationContext) LabelManager
	TypeManagerCreator     func(*GenerationContext) TypeManager
	GlobalManagerCreator   func(*GenerationContext) GlobalManager
}

// This structure is the most broad level of generation context.
// It contains information that is used in different parts of the compilation,
// and which are essential across the whole pipeline.
//
// It mainly contains information about the compilation target.
type GenerationContext struct {
	ManagerCreators

	// An instruction (definition) manager, which contains all instruction
	// definitions that are supported in the current architecture (ISA).
	//
	// When processing a new instruction from the source code, the compiler
	// talks with the instruction manager, retrieves the relevant instruction
	// definition, and uses it to farther process the instruction.
	Instructions InstructionManager

	// The size of a pointer type in the current target architecture.
	//
	// TODO: I'm not sure that we need this information in this step of the
	//   compilation. Can we compile without it, and leave it to the ISA?
	PointerSize *big.Int
}

func (ctx *GenerationContext) NewFileGenerationContext(
	source core.SourceContext,
) *FileGenerationContext {
	return &FileGenerationContext{
		GenerationContext: ctx,
		SourceContext:     source,
		Types:             ctx.TypeManagerCreator(ctx),
		Globals:           ctx.GlobalManagerCreator(ctx),
	}
}

type FileGenerationContext struct {
	*GenerationContext

	// The source code of the file that we are currently processing.
	core.SourceContext

	// A type manager that contains all types defined in the file,
	// and that the compiler can use to create new types when in processes
	// new type definitions.
	Types TypeManager

	// A type manager that contains all declared and defines globals in the file.
	// This includes function definitions and "extern" declarations.
	Globals GlobalManager
}

func (ctx *FileGenerationContext) NewFunctionGenerationContext() *FunctionGenerationContext {
	return &FunctionGenerationContext{
		FileGenerationContext: ctx,
		Registers:             ctx.RegisterManagerCreator(ctx),
		Labels:                ctx.LabelManagerCreator(ctx),
	}
}

type FunctionGenerationContext struct {
	*FileGenerationContext

	// A register manager, which contains all active registers in the function.
	// The compiler can query register information from the register manager
	// when it encounters registers while processing, and can create new register
	// information structures and pass them to the manager, which stores them.
	Registers RegisterManager

	// A label manager, which stores and manages all labels defined in a function.
	Labels LabelManager
}

type InstructionGenerationContext struct {
	*FunctionGenerationContext

	// A (partial) instruction info type for which we are currently working on
	// generating.
	InstructionInfo *InstructionInfo
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

type FunctionContextGenerator[NodeT parse.Node, InfoT any] interface {
	Generate(
		ctx *FunctionGenerationContext,
		node NodeT,
	) (info InfoT, results core.ResultList)
}

type InstructionContextGenerator[NodeT parse.Node, InfoT any] interface {
	Generate(
		ctx *InstructionGenerationContext,
		node NodeT,
	) (info InfoT, results core.ResultList)
}

// MARK: Utils

func ViewToSourceString(
	ctx *FileGenerationContext,
	view core.UnmanagedSourceView,
) string {
	return string(view.Raw(ctx.SourceContext))
}

func NodeToSourceString(
	ctx *FileGenerationContext,
	node parse.Node,
) string {
	return ViewToSourceString(ctx, node.View())
}
