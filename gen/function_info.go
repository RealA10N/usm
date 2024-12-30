package gen

type FunctionInfo struct {
	EntryBlock *BasicBlockInfo
	Registers  []*RegisterInfo
	Parameters []*RegisterInfo
	// TODO: add targets
}
