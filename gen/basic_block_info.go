package gen

type BasicBlockInfo struct {
	Labels       []*LabelInfo
	Instructions []*InstructionInfo

	ForwardEdges  []*BasicBlockInfo
	BackwardEdges []*BasicBlockInfo
}
