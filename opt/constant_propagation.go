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

// cpBlockSeparator is the sentinel pushed onto cpReachingConstants.registerPushes
// to mark the boundary between basic blocks in the DFS.
// Mirrors the blockSeparator sentinel used in ReachingDefinitionsSet (opt/ssa).
const cpBlockSeparator = ^uint(0)

// cpReachingConstants tracks, for each register, the constant value of its
// reaching definition at the current position in a DFS traversal.
type cpReachingConstants struct {
	// Per-register stacks of reaching constant values. Indices are assigned
	// lazily via registersToIndex the first time a register is defined.
	registerStacks []stack.Stack[*gen.ImmediateInfo]

	// Records which register indices were pushed in each block, with
	// cpBlockSeparator values marking block boundaries.
	registerPushes stack.Stack[uint]

	// Maps each *RegisterInfo to its index in registerStacks. Entries are
	// added on demand in define().
	registersToIndex map[*gen.RegisterInfo]uint

	// Set of reachable blocks (from entry). Used by define() to ignore
	// definitions in unreachable blocks when checking eligibility.
	reachable map[*gen.BasicBlockInfo]bool
}

func newCPReachingConstants(reachable map[*gen.BasicBlockInfo]bool) cpReachingConstants {
	return cpReachingConstants{
		registerStacks:   []stack.Stack[*gen.ImmediateInfo]{},
		registerPushes:   stack.New[uint](),
		registersToIndex: make(map[*gen.RegisterInfo]uint),
		reachable:        reachable,
	}
}

// get returns the constant value of the reaching definition for reg at the
// current DFS position, or nil if no constant reaching definition is known.
func (c *cpReachingConstants) get(reg *gen.RegisterInfo) *gen.ImmediateInfo {
	idx, ok := c.registersToIndex[reg]
	if !ok {
		return nil
	}
	stk := c.registerStacks[idx]
	if len(stk) == 0 {
		return nil
	}
	return stk[len(stk)-1]
}

// define records a new constant reaching definition for reg in currentBlock.
//
// A constant is only pushed when every reachable definition of reg is in
// currentBlock. Definitions in unreachable blocks are ignored. This covers
// two cases:
//   - Exactly one reachable definition, in currentBlock: it dominates all
//     uses, so the stack value is the true reaching value.
//   - Multiple reachable definitions, all in currentBlock: sequential
//     redefinitions within the same block; later ones shadow earlier ones.
//
// Any reachable definition in a different block disqualifies the register,
// because DFS ancestry does not imply domination at join points.
func (c *cpReachingConstants) define(
	reg *gen.RegisterInfo,
	imm *gen.ImmediateInfo,
	currentBlock *gen.BasicBlockInfo,
) {
	if imm == nil {
		return
	}

	for _, def := range reg.Definitions {
		if c.reachable[def.BasicBlockInfo] && def.BasicBlockInfo != currentBlock {
			return
		}
	}

	idx, ok := c.registersToIndex[reg]
	if !ok {
		idx = uint(len(c.registerStacks))
		c.registersToIndex[reg] = idx
		c.registerStacks = append(c.registerStacks, stack.New[*gen.ImmediateInfo]())
	}
	c.registerPushes.Push(idx)
	c.registerStacks[idx].Push(imm)
}

func (c *cpReachingConstants) pushBlock() {
	c.registerPushes.Push(cpBlockSeparator)
}

func (c *cpReachingConstants) popBlock() {
	for c.registerPushes.Top() != cpBlockSeparator {
		idx := c.registerPushes.Top()
		c.registerPushes.Pop()
		c.registerStacks[idx].Pop()
	}
	c.registerPushes.Pop() // pop the block separator itself
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
	for _, def := range cpInstruction.PropagateConstants(instruction) {
		reaching.define(def.Register, def.Immediate, instruction.BasicBlockInfo)
	}

	return core.ResultList{}
}

// ConstantPropagation replaces all uses of registers that are provably
// assigned a constant immediate value with that immediate directly.
//
// The pass performs a DFS from the entry block (block 0), visiting only
// reachable blocks. Unreachable blocks are skipped entirely.
//
// Within the DFS tree, define() only pushes a constant for a register when
// all of its definitions are confined to the current block. This keeps the
// analysis sound: registers defined in multiple blocks are skipped because
// DFS ancestry does not imply domination at join points.
//
// All instructions in the function must implement
// ConstantPropagationSupportedInstruction.
func ConstantPropagation(function *gen.FunctionInfo) core.ResultList {
	cfInfo := gen.NewFunctionControlFlowInfo(function)
	n := uint(len(cfInfo.BasicBlocks))
	dfs := cfInfo.ControlFlowGraph.Dfs(0)

	reachable := make(map[*gen.BasicBlockInfo]bool, n)
	for _, event := range dfs.Timeline {
		if event < n {
			reachable[cfInfo.BasicBlocks[event]] = true
		}
	}

	reaching := newCPReachingConstants(reachable)
	results := core.ResultList{}

	for _, event := range dfs.Timeline {
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
