package gen

import "alon.kr/x/usm/core"

type RegisterManager interface {
	GetRegister(name string) *RegisterInfo
	NewRegister(reg *RegisterInfo) core.Result
}
