package gen

import (
	"alon.kr/x/usm/core"
)

type BaseInstruction interface {
	// This method is usd by the USM engine to generate the internal control
	// flow graph representation.
	//
	// Should return a non-empty slice. If the instruction does not have any
	// consecutive steps in the function (for example, a return statement),
	// then a special dedicated return step should be returned.
	PossibleNextSteps() ([]StepInfo, core.ResultList)

	// Returns the string that represents the operator of the instruction.
	// For example, for the add instruction this method would return "ADD".
	String() string
}

// A basic instruction definition. This defines the logic that converts the
// generic, architecture / instruction set independent instruction AST nodes
// into a format instruction which is part of a specific instruction set.
type InstructionDefinition interface {
	// Build an instruction from the provided instruction information.
	BuildInstruction(info *InstructionInfo) (BaseInstruction, core.ResultList)

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
		ctx *FunctionGenerationContext,
		targets []*ReferencedTypeInfo,
		arguments []*ReferencedTypeInfo,
	) ([]ReferencedTypeInfo, core.ResultList)
}
