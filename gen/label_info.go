package gen

import "alon.kr/x/usm/core"

type LabelInfo[InstT BaseInstruction] struct {
	// The name of the label, as it appears in the source code.
	Name string

	BasicBlock *BasicBlockInfo[InstT]

	// A view of the label declaration in the source code.
	Declaration core.UnmanagedSourceView
}
