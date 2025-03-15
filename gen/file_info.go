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
