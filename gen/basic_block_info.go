package gen

import "slices"

type BasicBlockInfo struct {
	*FunctionInfo

	Label        *LabelInfo
	Instructions []*InstructionInfo

	ForwardEdges  []*BasicBlockInfo
	BackwardEdges []*BasicBlockInfo

	NextBlock *BasicBlockInfo
}

func (i *BasicBlockInfo) String() string {
	s := i.Label.String() + "\n"
	for _, instruction := range i.Instructions {
		s += "\t" + instruction.String() + "\n"
	}
	return s
}

func NewEmptyBasicBlockInfo(function *FunctionInfo) *BasicBlockInfo {
	return &BasicBlockInfo{
		FunctionInfo: function,
	}
}

// Returns the number of instructions in the basic block.
func (b *BasicBlockInfo) Size() int {
	return len(b.Instructions)
}

func (b *BasicBlockInfo) SetLabel(label *LabelInfo) {
	b.Label = label
	label.linkToBasicBlock(b)
}

func (b *BasicBlockInfo) AppendInstruction(instruction *InstructionInfo) {
	instruction.BasicBlockInfo = b
	b.Instructions = append(b.Instructions, instruction)
}

func (b *BasicBlockInfo) PrependInstruction(instruction *InstructionInfo) {
	// TODO: convert instructions to a linked list.
	instruction.BasicBlockInfo = b
	b.Instructions = append([]*InstructionInfo{instruction}, b.Instructions...)
}

func (b *BasicBlockInfo) RemoveInstruction(
	instruction *InstructionInfo,
) (ok bool) {
	instruction.BasicBlockInfo = nil

	instructionIndex := slices.Index(b.Instructions, instruction)
	if instructionIndex == -1 {
		return false
	}

	b.Instructions = slices.Delete(b.Instructions, instructionIndex, instructionIndex+1)
	return true
}

func (b *BasicBlockInfo) AppendBasicBlock(otherBlock *BasicBlockInfo) {
	otherBlock.NextBlock = b.NextBlock
	b.NextBlock = otherBlock
}

func (b *BasicBlockInfo) AppendForwardEdge(otherBlock *BasicBlockInfo) {
	b.ForwardEdges = append(b.ForwardEdges, otherBlock)
	otherBlock.BackwardEdges = append(otherBlock.BackwardEdges, b)
}
