package gen

type FunctionInfo struct {
	Name       string
	EntryBlock *BasicBlockInfo
	Registers  RegisterManager
	Labels     LabelManager
	Parameters []*RegisterInfo
	Targets    []ReferencedTypeInfo
}

func (f *FunctionInfo) CollectBasicBlocks() []*BasicBlockInfo {
	blocks := make([]*BasicBlockInfo, 0)
	for block := f.EntryBlock; block != nil; block = block.NextBlock {
		blocks = append(blocks, block)
	}

	return blocks
}

func (i *FunctionInfo) String() string {
	s := "func "
	for _, target := range i.Targets {
		s += target.String() + " "
	}

	s += i.Name + " "

	for _, param := range i.Parameters {
		// TODO: create a separate ParameterInfo type and just call String()
		// on it.
		s += param.Type.String() + " " + param.String() + " "
	}

	s += "{\n"

	for block := i.EntryBlock; block != nil; block = block.NextBlock {
		s += block.String()
	}

	s += "}"
	return s
}
