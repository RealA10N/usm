package ssa

import (
	"alon.kr/x/stack"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// A structure that holds the reaching definitions of registers while traversing
// a function in the SSA construction process.
//
// The traversal of basic blocks in the function SSA construction renaming phase
// is done in DFS order of the *Dominator tree*. Since After insertion of phi
// instructions in the first phase of the SSA construction, each register use
// is dominated exactly by one definition of the register (by the definition
// of the insertion points of phi instructions), in the dominator tree, given
// a basic block in which we use a register, the unique definition that reaches
// its use (disregarding other definitions in the same basic block) is the
// lowest ancestor of the basic block in the dominator tree that contains a
// definition of the register.
//
// We use the property described above in this data structure.
// Intuitively, during the traversal of the dominator tree, this data structure
// keeps a stack of all register definitions in the block from the currently
// traversed block to the root of the dominator tree (the entry node).
type ReachingDefinitionsSet struct {
	*FunctionSsaInfo

	// A mapping from a base register, originally defined in the function,
	// to it's reaching definitions. The top entry in each stack is the current
	// reaching definition for the currently active basic block.
	registerDefinitionStacks []stack.Stack[*gen.RegisterInfo]

	// A stack that holds the indices of the registers that were given a new
	// definition in the basic blocks from the entry node to the currently
	// processed node.
	//
	// Internally, each new beginning of a block in the path is marked by a
	// special `blockSeparator` value on the stack, and on top of it we push
	// the registers that are redefined in the basic block, until another
	// `blockSeparator` is pushed. In particular, when the stack is non-empty,
	// you can assume that the bottom of the stack contains the `blockSeparator`
	// value.
	registerDefinitionPushes stack.Stack[uint]
}

const blockSeparator = ^uint(0)

func NewReachingDefinitionsSet(function *FunctionSsaInfo) ReachingDefinitionsSet {
	registers := len(function.Registers)
	return ReachingDefinitionsSet{
		registerDefinitionStacks: make([]stack.Stack[*gen.RegisterInfo], registers),
		registerDefinitionPushes: stack.New[uint](),
	}
}

// Provided an original base register, this function returns the reaching
// renamed register definition that originates from the base register.
//
// This function should be called via an implementation of the
// RenameBasicBlock method, in which the reaching definition is the unique
// definition of the register that reaches (dominates) the entry of the
// basic block.
//
// First, the USM engine initializes the data structure with the reaching
// definitions to the entry of the basic block, NOT including the phi
// functions defined in the block.
// Then, the implementation of the ISA should update the data structure
// while renaming variables in a basic block, and when a new definition
// is reached inside a basic block (including a definition from a phi
// instruction), it the implementation should call the "rename register"
// method to get the new renamed register, and update the internal reaching
// definitions set.
func (s *ReachingDefinitionsSet) GetReachingDefinition(
	base *gen.RegisterInfo,
) (renamed *gen.RegisterInfo) {
	registerIndex := s.FunctionSsaInfo.RegistersToIndex[base]
	stack := s.registerDefinitionStacks[registerIndex]
	lastIndex := len(stack) - 1
	return stack[lastIndex]
}

// Update the reaching definition of a base register to the new renamed one.
// This method is called by the ISA implementation while processing a linear
// basic block, if a new definition of a register is found.
//
// Note that this DOES include definitions of registers in phi instructions.
func (s *ReachingDefinitionsSet) RenameDefinitionRegister(
	base *gen.RegisterInfo,
) *gen.RegisterInfo {
	baseIndex := s.FunctionSsaInfo.RegistersToIndex[base]
	s.registerDefinitionPushes.Push(baseIndex)

	renamed := s.SsaConstructionScheme.NewRenamedRegister(base)
	s.registerDefinitionStacks[baseIndex].Push(renamed)

	return renamed
}

func (s *ReachingDefinitionsSet) pushBlock() {
	s.registerDefinitionPushes.Push(blockSeparator)
}

func (s *ReachingDefinitionsSet) popBlock() {
	// We assume the bottom of the pushes stack contains the separator value,
	// so we do not check for the special case of an empty stack.
	for s.registerDefinitionPushes.Top() != blockSeparator {
		registerIndex := s.registerDefinitionPushes.Top()
		s.registerDefinitionStacks[registerIndex].Pop()
	}

	// Finally, we pop the separator value.
	s.registerDefinitionPushes.Pop()
}

type PhiInstruction interface {
	AddForwardingRegister(
		*gen.BasicBlockInfo,
		*gen.RegisterInfo,
	) core.ResultList
}

// This interface defines the process of renaming registers in the SSA
// construction process. Each ISA should implement this interface to provide
// support for the SSA construction process.
type SsaConstructionScheme interface {
	// Inserts a new phi instruction in the provided basic block, which is used
	// as a new definition to the provided basic block. This is called before
	// the register renaming procedure, and thus the provided register is
	// a register which is defined in the source code.
	NewPhiInstruction(
		*gen.BasicBlockInfo,
		*gen.RegisterInfo,
	) (PhiInstruction, core.ResultList)

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
