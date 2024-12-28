package gen

type FunctionInfo struct {
	EntryBlock *BasicBlockInfo
	Parameters []*RegisterInfo
	// TODO: add targets
}
