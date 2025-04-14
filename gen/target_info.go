package gen

import "alon.kr/x/usm/core"

type TargetInfo struct {
	Register    *RegisterInfo
	Declaration *core.UnmanagedSourceView
}

func NewTargetInfo(register *RegisterInfo) TargetInfo {
	return TargetInfo{
		Register:    register,
		Declaration: nil,
	}
}

func (i *TargetInfo) String() string {
	return i.Register.Type.String() + " " + i.Register.String()
}
