package gen

type BasicBlockInfo struct {
	Labels       []*LabelInfo
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
