package gen

type FunctionInfo struct {
	EntryBlock *BasicBlockInfo
	Registers  []*RegisterInfo
	Parameters []*RegisterInfo
	Targets    []ReferencedTypeInfo
}
