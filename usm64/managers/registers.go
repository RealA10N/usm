package managers

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

func (m *RegisterMap) NewRegister(reg *gen.RegisterInfo) core.Result {
	(*m)[reg.Name] = reg
	return nil
}

func (m *RegisterMap) GetAllRegisters() []*gen.RegisterInfo {
	registers := make([]*gen.RegisterInfo, 0, len(*m))
	for _, reg := range *m {
		registers = append(registers, reg)
	}
	return registers
}

func NewRegisterManager() gen.RegisterManager {
	return &RegisterMap{}
}
