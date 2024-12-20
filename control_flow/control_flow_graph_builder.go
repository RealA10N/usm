package control_flow

func getInstructionsForwardEdges[InstT SupportsControlFlow](
	instructions []InstT,
) [][]uint {
	forwardEdges := make([][]uint, len(instructions))
	for i, inst := range instructions {
		forwardEdges[i] = inst.PossibleNextInstructionIndices()
	}
	return forwardEdges
}

func backwardEdgesFromForwardEdges(forwardEdges [][]uint) [][]uint {
	backwardEdges := make([][]uint, len(forwardEdges))
	for from := range forwardEdges {
		for _, to := range forwardEdges[from] {
			backwardEdges[to] = append(backwardEdges[to], uint(from))
		}
	}
	return backwardEdges
}

type controlFlowGraphBuilder struct {
	ForwardEdges            [][]uint
	BackwardEdges           [][]uint
	BasicBlocks             []ControlFlowBasicBlock
	Visited                 []bool
	InstructionToBasicBlock []uint
}

func (b *controlFlowGraphBuilder) isLastInstructionInBasicBlock(instruction uint) bool {
	if len(b.ForwardEdges[instruction]) != 1 {
		return true
	}

	next := b.ForwardEdges[instruction][0]
	if len(b.BackwardEdges[next]) != 1 || b.Visited[next] {
		return true
	}

	return false
}

func (b *controlFlowGraphBuilder) addInstructionToBasicBlock(
	instruction uint,
	blockNumber uint,
) {
	b.Visited[instruction] = true
	b.InstructionToBasicBlock[instruction] = blockNumber
	b.BasicBlocks[blockNumber].NodeIndices = append(
		b.BasicBlocks[blockNumber].NodeIndices,
		instruction,
	)
}

func (b *controlFlowGraphBuilder) exploreBasicBlock(current uint) {
	if b.Visited[current] {
		return // already visited and processed in the past.
	}

	// current is the first instruction in a new basic block.
	// we now explore the basic block by just following forward edges,
	// until we reach an instruction that has more than 1 outgoing edge
	// or incoming edges (or zero?).

	blockNumber := uint(len(b.BasicBlocks))
	b.BasicBlocks = append(b.BasicBlocks, ControlFlowBasicBlock{
		NodeIndices:   []uint{},
		ForwardEdges:  []uint{},
		BackwardEdges: []uint{},
	})

	// traverse the current basic block while we can.

	for !b.isLastInstructionInBasicBlock(current) {
		b.addInstructionToBasicBlock(current, blockNumber)
		current = b.ForwardEdges[current][0]
	}

	// finish traversal: update last instruction in the basic block.
	b.addInstructionToBasicBlock(current, blockNumber)

	// explore following basic blocks recursively.
	for _, next := range b.ForwardEdges[current] {
		b.exploreBasicBlock(next)
		nextBlock := b.InstructionToBasicBlock[next]
		b.BasicBlocks[blockNumber].ForwardEdges = append(
			b.BasicBlocks[blockNumber].ForwardEdges,
			nextBlock,
		)
		b.BasicBlocks[nextBlock].BackwardEdges = append(
			b.BasicBlocks[nextBlock].BackwardEdges,
			blockNumber,
		)
	}
}

func NewControlFlowGraph(nodes []SupportsControlFlow) ControlFlowGraph {
	// A new basic blocks begins if the current instruction (first instruction
	// in the in the block), is:
	// (1) The first instruction in the function.
	// (2) There is more than 1 incoming edge to the instruction.
	//
	// Similarly, a basic block ends (or, splits into multiple basic blocks),
	// if the current instruction (last instruction in the block), is:
	// (1) The last instruction in the function.
	// (2) Has more than 1 outgoing edge from the instruction.
	//
	// The SupportsCFG interface gives as an easy API for forward (outgoing)
	// edges only. The easiest (and probably not the fastest!) implementation
	// is to first iterate over all instructions, and compute all backward
	// edges. Then, perform a DFS using the rules above to create basic blocks.

	forwardEdges := getInstructionsForwardEdges(nodes)
	backwardEdges := backwardEdgesFromForwardEdges(forwardEdges)

	builder := controlFlowGraphBuilder{
		ForwardEdges:            forwardEdges,
		BackwardEdges:           backwardEdges,
		BasicBlocks:             make([]ControlFlowBasicBlock, 0),
		Visited:                 make([]bool, len(nodes)),
		InstructionToBasicBlock: make([]uint, len(nodes)),
	}

	for i := range nodes {
		builder.exploreBasicBlock(uint(i))
	}

	return ControlFlowGraph{
		BasicBlocks: builder.BasicBlocks,
	}
}
