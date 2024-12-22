package graph

type controlFlowGraphBuilder struct {
	// The original graph
	*Graph

	// The new control flow graph that we are building.
	ControlFlowGraph

	// Visited[v] is true if v in (the original graph) it has already been
	// processed by the DFS construction algorithm. Initialized to false for all
	// vertices.
	Visited []bool
}

func (b *controlFlowGraphBuilder) hasExactlyOneForwardEdge(instruction uint) bool {
	return len(b.Graph.Nodes[instruction].ForwardEdges) == 1
}

func (b *controlFlowGraphBuilder) hasExactlyOneBackwardEdge(instruction uint) bool {
	return len(b.Graph.Nodes[instruction].BackwardEdges) == 1
}

// Checks if the provided instruction node can be "merged" with the next
// instruction node in the original graph, into a single basic block.
//
// Essentially, it checks that the instruction has only one forward edge, and
// that the next instruction has only one backwards edge, and that we have not
// yet visited the next instruction (to avoid infinite loops while processing
// unreachable cycles).
func (b *controlFlowGraphBuilder) isLastInstructionInBasicBlock(instruction uint) bool {
	if !b.hasExactlyOneForwardEdge(instruction) {
		return true
	}

	next := b.Graph.Nodes[instruction].ForwardEdges[0]
	if b.Visited[next] {
		return true
	}

	if !b.hasExactlyOneBackwardEdge(next) {
		return true
	}

	return false
}

func (b *controlFlowGraphBuilder) addInstructionToBasicBlock(
	node uint,
	block uint,
) {
	b.Visited[node] = true
	b.ControlFlowGraph.NodeToBasicBlock[node] = block
	b.ControlFlowGraph.BasicBlockToNodes[block] = append(
		b.ControlFlowGraph.BasicBlockToNodes[block],
		node,
	)
}

func (b *controlFlowGraphBuilder) createNewBasicBlock() uint {
	block := b.ControlFlowGraph.Size()
	b.ControlFlowGraph.BasicBlockToNodes = append(b.ControlFlowGraph.BasicBlockToNodes, nil)
	b.ControlFlowGraph.Nodes = append(b.ControlFlowGraph.Nodes, Node{})
	return block
}

func (b *controlFlowGraphBuilder) getNextInstruction(instruction uint) uint {
	return b.Graph.Nodes[instruction].ForwardEdges[0]
}

func (b *controlFlowGraphBuilder) exploreBasicBlock(current uint) {
	// A new basic blocks begins if the current instruction (first instruction
	// in the in the block), is:
	// (1) The first instruction in the function.
	// (2) There is more than 1 incoming edge to the instruction.
	//
	// Similarly, a basic block ends (or, splits into multiple basic blocks),
	// if the current instruction (last instruction in the block), is:
	// (1) The last instruction in the function.
	// (2) Has more than 1 outgoing edge from the instruction.

	if b.Visited[current] {
		return // already visited and processed in the past.
	}

	// current is the first instruction in a new basic block.
	// we now explore the basic block by just following forward edges,
	// until we reach an instruction that has more than 1 outgoing edge
	// or incoming edges (or zero?).

	block := b.createNewBasicBlock()

	// traverse the current basic block while we can.
	for !b.isLastInstructionInBasicBlock(current) {
		b.addInstructionToBasicBlock(current, block)
		current = b.getNextInstruction(current)
	}

	// finish traversal: update last instruction in the basic block.
	b.addInstructionToBasicBlock(current, block)

	// explore following basic blocks recursively.
	for _, next := range b.Graph.Nodes[current].ForwardEdges {
		b.exploreBasicBlock(next)
		nextBlock := b.ControlFlowGraph.NodeToBasicBlock[next]
		b.ControlFlowGraph.AddEdge(block, nextBlock)
	}
}

func newNodeToBasicBlock(size uint) []uint {
	nodeToBasicBlock := make([]uint, size)
	for i := uint(0); i < size; i++ {
		nodeToBasicBlock[i] = Unreachable
	}
	return nodeToBasicBlock
}

func newControlFlowGraphBuilder(g *Graph) controlFlowGraphBuilder {
	return controlFlowGraphBuilder{
		Graph: g,
		ControlFlowGraph: ControlFlowGraph{
			Graph:             NewEmptyGraph(0),
			BasicBlockToNodes: make([][]uint, 0),
			NodeToBasicBlock:  newNodeToBasicBlock(g.Size()),
		},
		Visited: make([]bool, len(g.Nodes)),
	}
}
