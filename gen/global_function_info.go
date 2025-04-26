package gen

import "alon.kr/x/usm/core"

type FunctionGlobalInfo struct {
	*FunctionInfo
}

func NewFunctionGlobalInfo(functionInfo *FunctionInfo) GlobalInfo {
	return &FunctionGlobalInfo{
		FunctionInfo: functionInfo,
	}
}

func (f *FunctionGlobalInfo) Name() string {
	return f.FunctionInfo.Name
}

func (f *FunctionGlobalInfo) Declaration() *core.UnmanagedSourceView {
	return f.FunctionInfo.Declaration
}
