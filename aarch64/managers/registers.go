package aarch64managers

import (
	"alon.kr/x/aarch64codegen/registers"
	"alon.kr/x/list"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

const numOfRegisters = registers.GPRegisterXZR + 1

// Currently, only 64 bit registers are supported in aarch64.
// There are 31 general purpose registers, named X0-X30, and one zero register,
// XZR. Although the zero register is not a general purpose register, it can be
// used as a general purpose register in instructions, so it is included here.
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
	return list.FromSingle(core.Result{
		{
			Type:     core.ErrorResult,
			Message:  "AArch64 does not allow defining new registers",
			Location: &reg.Declaration,
		},
	})
}

func (m *Aarch64RegisterManager) DeleteRegister(reg *gen.RegisterInfo) core.ResultList {
	return list.FromSingle(core.Result{
		{
			Type:     core.InternalErrorResult,
			Message:  "Trying to delete a register in AArch64",
			Location: &reg.Declaration,
		},
	})
}

func (m *Aarch64RegisterManager) Size() int {
	return len(m.registers)
}

func (m *Aarch64RegisterManager) GetAllRegisters() []*gen.RegisterInfo {
	registers := []*gen.RegisterInfo{}
	for _, register := range m.registers {
		registers = append(registers, register)
	}

	return registers
}

func NewRegisterManager(fileCtx *gen.FileGenerationContext) gen.RegisterManager {
	type64 := fileCtx.Types.GetType("$64")
	regs := [numOfRegisters]*gen.RegisterInfo{}

	for gpr := registers.GPRegister(0); gpr < numOfRegisters; gpr++ {
		info := gen.NewRegisterInfo(
			"%"+gpr.String(),
			gen.ReferencedTypeInfo{Base: type64},
		)

		regs[gpr] = info
	}

	return &Aarch64RegisterManager{
		registers: regs,
	}
}
