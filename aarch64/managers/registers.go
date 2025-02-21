package aarch64managers

import (
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

const numOfRegisters = 31

// Currently, only 64 bit registers are supported in aarch64.
// There are 31 general purpose registers, named x0-x30.
type Aarch64RegisterManager struct {
	// Registers are defined lazily, where X{i} is stored in register[i] if defined,
	// or is nil if it is not defined.
	registers [numOfRegisters]*gen.RegisterInfo
}

func (m *Aarch64RegisterManager) GetRegister(name string) *gen.RegisterInfo {
	gpr, ok := aarch64translation.RegisterNameToAarch64GPRegister(name)
	if !ok {
		return nil
	}

	return m.registers[gpr]
}

func (m *Aarch64RegisterManager) NewRegister(reg *gen.RegisterInfo) core.ResultList {
	gpr, results := aarch64translation.RegisterToAarch64GPRegister(reg)
	if !results.IsEmpty() {
		return results
	}

	m.registers[gpr] = reg
	return core.ResultList{}
}

func (m *Aarch64RegisterManager) DeleteRegister(reg *gen.RegisterInfo) core.ResultList {
	gpr, results := aarch64translation.RegisterToAarch64GPRegister(reg)
	if !results.IsEmpty() {
		return results
	}

	m.registers[gpr] = nil
	return core.ResultList{}
}

func (m *Aarch64RegisterManager) Size() uint {
	size := uint(0)
	for _, register := range m.registers {
		if register != nil {
			size++
		}
	}

	return size
}

func (m *Aarch64RegisterManager) GetAllRegisters() []*gen.RegisterInfo {
	registers := []*gen.RegisterInfo{}
	for _, register := range m.registers {
		registers = append(registers, register)
	}

	return registers
}

func NewRegisterManager() gen.RegisterManager {
	return &Aarch64RegisterManager{}
}
