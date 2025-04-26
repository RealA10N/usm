package gen

import (
	"alon.kr/x/usm/core"
)

type GlobalMap map[string]GlobalInfo

func NewGlobalMap(*GenerationContext) GlobalManager {
	return &GlobalMap(make(map[string]GlobalInfo))
}

func (m *GlobalMap) GetGlobal(name string) GlobalInfo {
	return (*m)[name]
}

func (m *GlobalMap) NewGlobal(global GlobalInfo) core.ResultList {
	name := global.Name()
	(*m)[name] = global
	return core.ResultList{}
}

func (m *GlobalMap) GetAllGlobals() []GlobalInfo {
	globals := make([]GlobalInfo, 0, len(*m))
	for _, global := range *m {
		globals = append(globals, global)
	}
	return globals
}
