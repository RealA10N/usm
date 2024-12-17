package control_flow

type SupportsControlFlow interface {
	// Returns a list of all instruction indices in the function that execution
	// could arrive to after the execution of this instruction.
	PossibleNextInstructionIndices() []uint
}

type ControlFlowBasicBlock struct {
	InstructionIndices []uint
	ForwardEdges       []uint
	BackwardEdges      []uint
}

type ControlFlowGraph[InstT SupportsControlFlow] struct {
	Instructions []InstT
	BasicBlocks  []ControlFlowBasicBlock
}

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
	ForwardEdges  [][]uint
	BackwardEdges [][]uint
	BasicBlocks   []ControlFlowBasicBlock
	Visited       []bool
}

func (b *controlFlowGraphBuilder) isEndOfBasicBlock(instruction uint) bool {
	if len(b.ForwardEdges[instruction]) != 1 {
		return true
	}

	next := b.ForwardEdges[instruction][0]
	if len(b.BackwardEdges[next]) != 1 {
		return true
	}

	return false
}

func (b *controlFlowGraphBuilder) exploreBasicBlock(current uint) {
	if b.Visited[current] {
		return // already visited and processed in the past.
	}

	b.Visited[current] = true

	// current is the first instruction in a new basic block.
	// we now explore the basic block by just following forward edges,
	// until we reach an instruction that has more than 1 outgoing edge
	// or incoming edges (or zero?).

	block := ControlFlowBasicBlock{
		InstructionIndices: []uint{current},
		BackwardEdges:      b.BackwardEdges[current],
	}

	// traverse the current basic block while we can.

	for !b.isEndOfBasicBlock(current) {
		next := b.ForwardEdges[current][0]
		b.Visited[next] = true
		block.InstructionIndices = append(block.InstructionIndices, current)
		current = next
	}

	// finished traversing the basic block: update metadata

	block.ForwardEdges = b.ForwardEdges[current]
	b.BasicBlocks = append(b.BasicBlocks, block)

	// explore following basic blocks recursively.

	for _, next := range block.ForwardEdges {
		b.exploreBasicBlock(next)
	}
}

func NewControlFlowGraph[InstT SupportsControlFlow](
	instructions []InstT,
) ControlFlowGraph[InstT] {
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

	forwardEdges := getInstructionsForwardEdges(instructions)
	backwardEdges := backwardEdgesFromForwardEdges(forwardEdges)

	builder := controlFlowGraphBuilder{
		ForwardEdges:  forwardEdges,
		BackwardEdges: backwardEdges,
		BasicBlocks:   make([]ControlFlowBasicBlock, 0),
		Visited:       make([]bool, len(instructions)),
	}

	for i := range instructions {
		builder.exploreBasicBlock(uint(i))
	}

	return ControlFlowGraph[InstT]{
		Instructions: instructions,
		BasicBlocks:  builder.BasicBlocks,
	}
}
