package gen

type BasicBlockInfo struct {
	Instructions []*InstructionInfo

	ForwardEdges  []*BasicBlockInfo
	BackwardEdges []*BasicBlockInfo

	// All basic blocks in a function have a defined ordering between them.
	// The initial ordering that the USM engine produces is the order in which
	// the basic blocks appear in the source code.
	// the `NextBlock` field points to the next block that follows this block
	// in the ordering, or nil if this is the last basic block in the function.
	NextBlock *BasicBlockInfo
}

func (i *BasicBlockInfo) String() string {
	s := ""
	for _, instruction := range i.Instructions {
		s += "\t" + instruction.String() + "\n"
	}
	return s
}

func NewEmptyBasicBlockInfo() *BasicBlockInfo {
	return &BasicBlockInfo{
		Instructions: []*InstructionInfo{},

		ForwardEdges:  []*BasicBlockInfo{},
		BackwardEdges: []*BasicBlockInfo{},

		NextBlock: nil,
	}
}

func (b *BasicBlockInfo) AppendInstruction(instruction *InstructionInfo) {
	b.Instructions = append(b.Instructions, instruction)
	instruction.LinkToBasicBlock(b)
}
