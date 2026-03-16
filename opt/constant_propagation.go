package opt

import (
	"alon.kr/x/list"
	"alon.kr/x/stack"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// ConstantDefinition pairs a target register with the constant immediate
// value that an instruction always assigns to it.
type ConstantDefinition struct {
	Register  *gen.RegisterInfo
	Immediate *gen.ImmediateInfo
}

// ConstantPropagationSupportedInstruction is implemented by ISA instructions
// that support the constant propagation optimization pass.
type ConstantPropagationSupportedInstruction interface {
	gen.InstructionDefinition

	// PropagateConstants returns (register, constant) pairs for each target
	// that this instruction provably assigns a constant immediate to.
	// Returns an empty slice if this instruction defines no constants.
	// Example: "$64 %x = $64 #5" → [{Register: %x, Immediate: $64 #5}]
	PropagateConstants(info *gen.InstructionInfo) []ConstantDefinition
}

// PropagatesNoConstants can be embedded in instruction definitions that never
// assign a constant immediate to any of their targets.
type PropagatesNoConstants struct{}

func (PropagatesNoConstants) PropagateConstants(*gen.InstructionInfo) []ConstantDefinition {
	return nil
}

func newCPNotSupportedError(instruction *gen.InstructionInfo) core.ResultList {
	return list.FromSingle(core.Result{{
		Type:     core.InternalErrorResult,
		Message:  "Instruction does not support constant propagation",
		Location: instruction.Declaration,
	}})
}

// cpReachingConstants tracks, for each register, the constant value of its
// reaching definition at the current position in a DFS forest traversal.
//
// The design mirrors ReachingDefinitionsSet in opt/ssa: a per-register stack
// holds the sequence of constant values introduced along the path from the
// DFS tree root to the current block, and pushBlock/popBlock scope those
// values to the block's DFS subtree.
type cpReachingConstants struct {
	// Per-register stacks of reaching constant values, created on demand.
	registerStacks map[*gen.RegisterInfo]*stack.Stack[*gen.ImmediateInfo]

	// Records which registers were pushed in each block, with nil entries
	// marking block boundaries. Mirrors
	// ReachingDefinitionsSet.registerDefinitionPushes in opt/ssa.
	registerPushes stack.Stack[*gen.RegisterInfo]

	// Maps each basic block to the index of its DFS tree root (its component
	// ID). Used by define() to determine whether a definition belongs to the
	// same connected component as the block currently being processed.
	blockComponent map[*gen.BasicBlockInfo]uint
}

func newCPReachingConstants(
	blockComponent map[*gen.BasicBlockInfo]uint,
) cpReachingConstants {
	return cpReachingConstants{
		registerStacks: make(map[*gen.RegisterInfo]*stack.Stack[*gen.ImmediateInfo]),
		registerPushes: stack.New[*gen.RegisterInfo](),
		blockComponent: blockComponent,
	}
}

// get returns the constant value of the reaching definition for reg at the
// current DFS position, or nil if no constant reaching definition is known.
func (c *cpReachingConstants) get(reg *gen.RegisterInfo) *gen.ImmediateInfo {
	stk, ok := c.registerStacks[reg]
	if !ok || len(*stk) == 0 {
		return nil
	}
	return stk.Top()
}

// define records a new constant reaching definition for reg in currentBlock.
//
// A constant is pushed onto reg's stack only when every definition of reg in
// the same connected component as currentBlock is in currentBlock itself. This
// covers two cases uniformly:
//   - Exactly one in-component definition, in currentBlock: it dominates all
//     uses in well-formed code, so the stack value is always the true reaching
//     value.
//   - Multiple in-component definitions, all in currentBlock: they are
//     sequential redefinitions within the same block. Each define call pushes
//     one value; later definitions shadow earlier ones on the stack, and all
//     are popped together when the block is exited.
//
// Any in-component definition outside currentBlock disqualifies the register:
// DFS ancestry does not imply domination at join points, so propagating such a
// value would be unsound. Definitions in other components (dead code, isolated
// loops) do not count and do not disqualify the register.
func (c *cpReachingConstants) define(
	reg *gen.RegisterInfo,
	imm *gen.ImmediateInfo,
	currentBlock *gen.BasicBlockInfo,
) {
	if imm == nil {
		return
	}

	currentComponent := c.blockComponent[currentBlock]
	for _, def := range reg.Definitions {
		if c.blockComponent[def.BasicBlockInfo] == currentComponent && def.BasicBlockInfo != currentBlock {
			return
		}
	}

	stk, ok := c.registerStacks[reg]
	if !ok {
		newStk := stack.New[*gen.ImmediateInfo]()
		stk = &newStk
		c.registerStacks[reg] = stk
	}
	stk.Push(imm)
	c.registerPushes.Push(reg)
}

func (c *cpReachingConstants) pushBlock() {
	c.registerPushes.Push(nil)
}

func (c *cpReachingConstants) popBlock() {
	for c.registerPushes.Top() != nil {
		reg := c.registerPushes.Top()
		c.registerPushes.Pop()
		c.registerStacks[reg].Pop()
	}
	c.registerPushes.Pop() // pop the nil block separator itself
}

// cpProcessInstruction performs constant propagation for one instruction:
// it substitutes any register arguments whose reaching definition is a known
// constant, then records new reaching constants for targets whose value is
// now known after the substitution.
func cpProcessInstruction(
	instruction *gen.InstructionInfo,
	reaching *cpReachingConstants,
) core.ResultList {
	cpInstruction, ok := instruction.Definition.(ConstantPropagationSupportedInstruction)
	if !ok {
		return newCPNotSupportedError(instruction)
	}

	// Substitute arguments whose reaching definition is a known constant.
	for i, arg := range instruction.Arguments {
		regArg, ok := arg.(*gen.RegisterArgumentInfo)
		if !ok {
			continue
		}
		if imm := reaching.get(regArg.Register); imm != nil {
			instruction.SubstituteArgument(i, imm)
		}
	}

	// After substitution, record the reaching constant for each target.
	// PropagateConstants is called after substitution so that instructions
	// whose arguments just became constants (e.g. a move whose source was
	// just replaced) can propagate the constant to their own targets.
	for _, def := range cpInstruction.PropagateConstants(instruction) {
		reaching.define(def.Register, def.Immediate, instruction.BasicBlockInfo)
	}

	return core.ResultList{}
}

// ConstantPropagation replaces all uses of registers that are provably
// assigned a constant immediate value with that immediate directly.
//
// The pass uses DfsForest to traverse all basic blocks in a single unified
// DFS, covering every connected component of the CFG including unreachable
// blocks (isolated loops, dead code). Within each DFS tree the
// ancestor-before-descendant property ensures that any definition dominating
// a use is on the DFS stack when that use is processed.
//
// define() only pushes a constant for a register when all of its definitions
// in the same connected component are confined to the current block. This
// keeps the analysis sound for both SSA and non-SSA code: registers with
// definitions spanning multiple blocks in the same component are skipped
// because DFS ancestry does not imply domination at join points.
//
// All instructions in the function must implement
// ConstantPropagationSupportedInstruction.
func ConstantPropagation(function *gen.FunctionInfo) core.ResultList {
	cfInfo := gen.NewFunctionControlFlowInfo(function)
	n := uint(len(cfInfo.BasicBlocks))
	forest := cfInfo.ControlFlowGraph.DfsForest()

	// Compute the component ID for each block: the index of its DFS tree root.
	// Since parents always have a lower preorder than their children in a DFS
	// tree, iterating in preorder lets us propagate component IDs in one pass.
	componentOf := make([]uint, n)
	for preorder := uint(0); preorder < n; preorder++ {
		node := forest.PreOrderReversed[preorder]
		if forest.Parent[node] == node {
			componentOf[node] = node
		} else {
			componentOf[node] = componentOf[forest.Parent[node]]
		}
	}

	blockComponent := make(map[*gen.BasicBlockInfo]uint, n)
	for i, block := range cfInfo.BasicBlocks {
		blockComponent[block] = componentOf[uint(i)]
	}

	reaching := newCPReachingConstants(blockComponent)
	results := core.ResultList{}

	for _, event := range forest.Timeline {
		if event >= n {
			reaching.popBlock()
			continue
		}

		reaching.pushBlock()
		block := cfInfo.BasicBlocks[event]
		for _, instruction := range block.Instructions {
			curResults := cpProcessInstruction(instruction, &reaching)
			results.Extend(&curResults)
		}
	}

	return results
}
