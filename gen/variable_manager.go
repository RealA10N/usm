package gen

import "alon.kr/x/usm/core"

type VariableManager interface {
	// Returns variable information by its name, or nil if no variable with
	// that name is managed by this manager.
	GetVariable(name string) *VariableInfo

	// Adds a new variable to the manager.
	// The implementation can assume the variable is not already registered
	// (i.e. GetVariable(variable.Name) == nil).
	NewVariable(variable *VariableInfo) core.ResultList

	// Returns all currently registered variables, in declaration order.
	GetAllVariables() []*VariableInfo
}
