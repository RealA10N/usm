package gen

import "alon.kr/x/usm/core"

type LabelInfo struct {
	// The name of the label, as it appears in the source code.
	Name string

	// The basic block to which the instruction points.
	BasicBlock *BasicBlockInfo

	// A view of the label declaration in the source code.
	Declaration core.UnmanagedSourceView
}

func (i *LabelInfo) LinkToBasicBlock(basicBlock *BasicBlockInfo) {
	i.BasicBlock = basicBlock
}

func (i *LabelInfo) String() string {
	return i.Name
}
