package control_flow

type SupportsControlFlow interface {
	// Returns a list of all instruction indices in the function that execution
	// could arrive to after the execution of this instruction.
	PossibleNextInstructionIndices() []uint
}

type ControlFlowBasicBlock struct {
	InstructionIndices []uint
	ForwardEdges       []uint
}

// BasicBlocks[entryNode] is the implicit entry block.
const entryNode = 0

type ControlFlowGraph[InstT SupportsControlFlow] struct {
	Instructions []InstT
	BasicBlocks  []ControlFlowBasicBlock
}

// Returns the list of instruction indices in their pre-order traversal order.
func (g *ControlFlowGraph[InstT]) PreOrderDfs() []uint {
	builder := DfsBuilder[InstT]{
		ControlFlowGraph: g,
		Visited:          make([]bool, len(g.BasicBlocks)),
		Order:            make([]uint, len(g.BasicBlocks)),
		NextTime:         0,
	}

	builder.preOrderDfs(0)
	return builder.Order
}
