package usm64managers

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type RegisterMap map[string]*gen.RegisterInfo

func (m *RegisterMap) GetRegister(name string) *gen.RegisterInfo {
	val, ok := (*m)[name]
	if !ok {
		return nil
	}
	return val
}

func (m *RegisterMap) NewRegister(reg *gen.RegisterInfo) core.ResultList {
	(*m)[reg.Name] = reg
	return core.ResultList{}
}

func (m *RegisterMap) DeleteRegister(reg *gen.RegisterInfo) core.ResultList {
	delete(*m, reg.Name)
	return core.ResultList{}
}

func (m *RegisterMap) Size() int {
	return len(*m)
}

func (m *RegisterMap) GetAllRegisters() []*gen.RegisterInfo {
	registers := make([]*gen.RegisterInfo, 0, len(*m))
	for _, reg := range *m {
		registers = append(registers, reg)
	}
	return registers
}

func NewRegisterManager(*gen.FileGenerationContext) gen.RegisterManager {
	return &RegisterMap{}
}
