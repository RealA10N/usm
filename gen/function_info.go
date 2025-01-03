package gen

type FunctionInfo struct {
	Name       string
	EntryBlock *BasicBlockInfo
	Registers  []*RegisterInfo
	Parameters []*RegisterInfo
	Targets    []ReferencedTypeInfo
}

func (i *FunctionInfo) String() string {
	s := "func "
	for _, target := range i.Targets {
		s += target.String() + " "
	}

	s += i.Name + " "

	for _, param := range i.Parameters {
		s += param.String() + " "
	}

	s += "{\n"

	block := i.EntryBlock
	for block != nil {
		s += block.String()
		block = block.NextBlock
	}

	s += "}"
	return s
}
