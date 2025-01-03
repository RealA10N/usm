package gen

type BasicBlockInfo struct {
	*FunctionInfo

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
		s += instruction.String() + "\n"
	}
	return s
}

// Get a single label instance that represents the basic block, if it exists.
// If the basic block has multiple labels, the function will return one of
// them arbitrary.
// If the block does not have labels pointing to it, returns nil.
func (i *BasicBlockInfo) GetRepresentingLabel() *LabelInfo {
	firstInstruction := i.Instructions[0]
	labels := firstInstruction.Labels

	if len(labels) < 1 {
		// TODO: in this case, we should generate a new label, using some sort of
		// "LabelManager".
		return nil
	}

	return labels[0]
}

func NewEmptyBasicBlockInfo(function *FunctionInfo) *BasicBlockInfo {
	return &BasicBlockInfo{
		FunctionInfo: function,
	}
}

func (b *BasicBlockInfo) AppendInstruction(instruction *InstructionInfo) {
	b.Instructions = append(b.Instructions, instruction)
	instruction.LinkToBasicBlock(b)
}

func (b *BasicBlockInfo) PrependInstruction(instruction *InstructionInfo) {
	// TODO: move labels to this instruction instead of the second one?
	// TODO: convert instructions to a linked list.
	b.Instructions = append([]*InstructionInfo{instruction}, b.Instructions...)
	instruction.LinkToBasicBlock(b)
}
