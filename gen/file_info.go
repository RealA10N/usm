package gen

type FileInfo struct {
	Functions []*FunctionInfo
}

func (i *FileInfo) String() string {
	s := ""
	for _, function := range i.Functions {
		s += function.String() + "\n"
	}
	return s
}
