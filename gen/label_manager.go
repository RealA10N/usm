package gen

import "alon.kr/x/usm/core"

type LabelManager interface {
	// Get a information about a label from its name.

	// Create a new label, provided it's information.
	// This should register the label in the label manager database as a valid
	// label, and future queries should be able to find it.
	GetLabel(name string) *LabelInfo

	// Generate a new label that will be used with the provided block.
	// This should create a new label with a unique name, but it SHOULDN'T add
	// it to the labels database (does not call NewLabel).
	GenerateLabel(block *BasicBlockInfo) *LabelInfo

	// Create a new label, provided it's information.
	// This should register the label in the label manager database as a valid
	// label, and future queries should be able to find it.
	NewLabel(info *LabelInfo) core.Result

	// Generate a new label that will be used with the provided block.
	// This should create a new label with a unique name, but it SHOULDN'T add
	// it to the labels database (does not call NewLabel).
	GenerateLabel(block *BasicBlockInfo) *LabelInfo
}
