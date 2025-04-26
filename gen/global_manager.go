package gen

import "alon.kr/x/usm/core"

type GlobalManager interface {
	GetGlobal(name string) GlobalInfo
	NewGlobal(GlobalInfo) core.ResultList
	GetAllGlobals() []GlobalInfo
}
