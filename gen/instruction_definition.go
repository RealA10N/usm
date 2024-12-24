package gen

import "alon.kr/x/usm/core"

// TODO: add basic interface methods for instruction.
type BaseInstruction interface{}

// A basic instruction definition. This defines the logic that converts the
// generic, architecture / instruction set independent instruction AST nodes
// into a format instruction which is part of a specific instruction set.
type InstructionDefinition[InstT BaseInstruction] interface {
	// Build an instruction from the provided targets and arguments.
	BuildInstruction(
		targets []*RegisterArgumentInfo,
		arguments []ArgumentInfo,
	) (InstT, core.ResultList)

	// Provided a list a list of types that correspond to argument types,
	// and a (possibly partial) list of target types, return a complete list
	// of target types which is implicitly inferred from the argument types,
	// and possibly the explicit target types, or an error if the target types
	// can not be inferred.
	//
	// On success, the length of the returned type slice should be equal to the
	// provided (partial) targets length. The non nil provided target types
	// should not be modified.
	//
	// TODO: perhaps we should not pass the bare generation context to the "public"
	// instruction set definition API, and should wrap it with a limited interface.
	InferTargetTypes(
		ctx *FunctionGenerationContext[InstT],
		targets []*ReferencedTypeInfo,
		arguments []*ReferencedTypeInfo,
	) ([]ReferencedTypeInfo, core.ResultList)
}
