package gen

import "alon.kr/x/usm/core"

type LabelManager interface {
	// Get a information about a label from its name.
	GetLabel(name string) *LabelInfo

	// Create a new label, provided it's information.
	// This should register the label in the label manager database as a valid
	// label, and future queries should be able to find it.
	NewLabel(info *LabelInfo) core.ResultList

	// Generate a new label that will be used to identify an unlabelled block.
	// This should create a new label with a unique name, but it SHOULDN'T add
	// it to the labels database (does not call NewLabel).
	GenerateLabel() *LabelInfo
}
