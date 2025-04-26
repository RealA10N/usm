package gen

import "alon.kr/x/usm/core"

type GlobalArgumentInfo struct {
	GlobalInfo

	declaration *core.UnmanagedSourceView
}

func (g *GlobalArgumentInfo) Declaration() *core.UnmanagedSourceView {
	return g.declaration
}

func (g *GlobalArgumentInfo) String() string {
	return g.GlobalInfo.Name()
}
