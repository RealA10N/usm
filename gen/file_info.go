package gen

import "alon.kr/x/usm/core"

type FileInfo struct {
	Functions map[string]*FunctionInfo
}

func NewFileInfo() *FileInfo {
	return &FileInfo{
		Functions: make(map[string]*FunctionInfo),
	}
}

func (i *FileInfo) String() string {
	s := ""

	if len(i.Functions) == 0 {
		return s
	}

	functions := make([]*FunctionInfo, 0, len(i.Functions))
	for _, function := range i.Functions {
		functions = append(functions, function)
	}

	for _, f := range functions[:len(functions)-1] {
		s += f.String() + "\n"
	}

	lastFunction := functions[len(functions)-1]
	s += lastFunction.String()

	return s
}

// GetFunction returns the function with the given name, or nil if it does not
// exist.
func (i *FileInfo) GetFunction(name string) *FunctionInfo {
	info, ok := i.Functions[name]
	if !ok {
		return nil
	}

	return info
}

func (i *FileInfo) AppendFunction(function *FunctionInfo) {
	oldFunction, ok := i.Functions[function.Name]
	if ok {
		oldFunction.FileInfo = nil
	}

	function.FileInfo = i
	i.Functions[function.Name] = function
}

func (i *FileInfo) Validate() core.ResultList {
	results := core.ResultList{}

	for _, function := range i.Functions {
		curResults := function.Validate()
		results.Extend(&curResults)
	}

	return results
}
