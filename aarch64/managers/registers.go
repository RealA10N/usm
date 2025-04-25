package aarch64managers

import (
	"alon.kr/x/faststringmap"
	"alon.kr/x/list"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

// Currently, only 64 bit registers are supported in aarch64.
// There are 31 general purpose registers, named x0-x30, one zero register named
// xzr, and one stack pointer register named sp.
type Aarch64RegisterManager struct {
	faststringmap.Map[*gen.RegisterInfo]
}

func (m *Aarch64RegisterManager) GetRegister(name string) *gen.RegisterInfo {
	reg, ok := m.LookupString(name)
	if !ok {
		return nil
	}

	return reg
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
	return len(aarch64translation.AllRegisterNames)
}

func (m *Aarch64RegisterManager) GetAllRegisters() []*gen.RegisterInfo {
	registers := make(
		[]*gen.RegisterInfo,
		0,
		len(aarch64translation.AllRegisterNames),
	)

	for _, name := range aarch64translation.AllRegisterNames {
		reg, ok := m.LookupString(name)
		if ok { // Should always be true
			registers = append(registers, reg)
		}
	}

	return registers
}

func NewRegisterManager(fileCtx *gen.FileGenerationContext) gen.RegisterManager {
	type64 := fileCtx.Types.GetType("$64")

	numOfRegisters := len(aarch64translation.AllRegisterNames)
	entries := make([]faststringmap.MapEntry[*gen.RegisterInfo], 0, numOfRegisters)

	for _, name := range aarch64translation.AllRegisterNames {
		info := gen.NewRegisterInfo(
			name,
			gen.ReferencedTypeInfo{Base: type64},
		)

		entry := faststringmap.MapEntry[*gen.RegisterInfo]{
			Key:   name,
			Value: info,
		}

		entries = append(entries, entry)
	}

	return &Aarch64RegisterManager{
		Map: faststringmap.NewMap[*gen.RegisterInfo](entries),
	}
}
