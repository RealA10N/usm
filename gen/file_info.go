package gen

type FileInfo struct {
	Functions []*FunctionInfo
}

func (i *FileInfo) String() string {
	s := ""

	if len(i.Functions) == 0 {
		return s
	}

	for _, function := range i.Functions[:len(i.Functions)-1] {
		s += function.String() + "\n"
	}

	s += i.Functions[len(i.Functions)-1].String()

	return s
}

// GetFunction returns the function with the given name, or nil if it does not
// exist.
func (i *FileInfo) GetFunction(name string) *FunctionInfo {
	for _, function := range i.Functions {
		if function.Name == name {
			return function
		}
	}
	return nil
}
