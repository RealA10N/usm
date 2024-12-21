package control_flow

type SupportsControlFlow interface {
	// Returns a list of all instruction indices in the function that execution
	// could arrive to after the execution of this instruction.
	PossibleNextInstructionIndices() []uint
}

type ControlFlowBasicBlock struct {
	NodeIndices   []uint
	ForwardEdges  []uint
	BackwardEdges []uint
}

// BasicBlocks[CfgEntryBlock] is the implicit entry block.
const CfgEntryBlock = 0

type ControlFlowGraph struct {
	BasicBlocks []ControlFlowBasicBlock
}

func (g *ControlFlowGraph) Size() uint {
	return uint(len(g.BasicBlocks))
}

// Returns the list of instruction indices in their pre-order traversal order.
func (g *ControlFlowGraph) Dfs(root uint) DfsResult {
	builder := newDfsBuilder(g)
	builder.dfs(root, root)
	return builder.toDfsResult()
}

func (g *ControlFlowGraph) DominatorTree() DominatorTree {
	builder := newDominatorTreeBuilder(g)
	immDom := builder.LengauerTarjan()

	return DominatorTree{
		ImmDom: immDom,
	}
}
