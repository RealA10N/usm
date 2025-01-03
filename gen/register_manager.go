package gen

import "alon.kr/x/usm/core"

type RegisterManager interface {
	// Returns register information by its name, or nil if the register is
	// not managed by this register manager.
	GetRegister(name string) *RegisterInfo

	// Adds a new register to the register manager.
	// The implementation can assume that the register is not already registered
	// (i.e. GetRegister(reg.Name) == nil).
	NewRegister(reg *RegisterInfo) core.ResultList

	// Deletes a register from the register manager.
	// The method should succeed even if the register is not found.
	DeleteRegister(reg *RegisterInfo) core.ResultList

	// Returns the number of currently registered registers.
	Size() uint

	// Returns all of the currently registered registers, as a slice.
	GetAllRegisters() []*RegisterInfo
}
