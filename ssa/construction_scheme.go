package ssa

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type ReachingDefinitionsSet interface {
	// Provided an original base register, this function returns the reaching
	// renamed register definition that originates from the base register.
	//
	// This function should be called via an implementation of the
	// RenameBasicBlock method, in which the reaching definition is the unique
	// definition of the register that reaches (dominates) the entry of the
	// basic block.

	// Note that this does not include definitions of the register in the basic
	// block itself, and it is up to the implementation to track.
	// This DOES NOT include definitions of the register in inserted phi
	// functions, including phi functions in the entry this basic block.
	// Such phi instructions should be first handled by the
	// RenameDefinitionRegister method, as if it was a regular register
	// definition.
	GetReachingDefinition(base *gen.RegisterInfo) (renamed *gen.RegisterInfo)

	// Update the reaching definition of a base register to the new renamed one.
	// This method is called by the ISA implementation while processing a linear
	// basic block, if a new definition of a register is found.
	//
	// Note that this DOES include definitions of registers in phi instructions.
	RenameDefinitionRegister(base *gen.RegisterInfo) (renamed *gen.RegisterInfo)
}

type PhiInstruction interface {
	AddForwardingRegister(*gen.BasicBlockInfo, *gen.RegisterInfo)
}

// This interface defines the process of renaming registers in the SSA
// construction process. Each ISA should implement this interface to provide
// support for the SSA construction process.
type SsaConstructionScheme interface {
	// Inserts a new phi instruction in the provided basic block, which is used
	// as a new definition to the provided basic block. This is called before
	// the register renaming procedure, and thus the provided register is
	// a register which is defined in the source code.
	NewPhiInstruction(*gen.BasicBlockInfo, *gen.RegisterInfo) PhiInstruction

	// Creates a new unique register that is used as a renaming of the provided
	// register in the construction of the SSA form.
	//
	// It is on the implementation to provide a unique register creation scheme,
	// which can create multiple new registers from the same base register.
	// The implementation must not return the same new register for the same
	// base register on different called, i.e. it should create a new register
	// every time it is called.
	//
	// Registers should not be renamed by the ISA implementation, and the
	// the USM engine does the renaming provided this interface.
	NewRenamedRegister(base *gen.RegisterInfo) (renamed *gen.RegisterInfo)

	// Provided a basic block and a set of the reaching definitions to that
	// basic block, the implementation should rename all registers that are
	// used and defined in the basic block.
	//
	// The caller does not grantee the order of basic blocks in which the calls
	// to this method are made.
	RenameBasicBlock(*gen.BasicBlockInfo, ReachingDefinitionsSet) core.ResultList
}
