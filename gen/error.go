package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
)

func NewUndefinedTypeResult(location core.UnmanagedSourceView) core.Result {
	return core.Result{{
		Type:     core.ErrorResult,
		Message:  "Undefined type",
		Location: &location,
	}}
}

func NewRegisterTypeMismatchResult(
	NewDeclaration core.UnmanagedSourceView,
	FirstDeclaration core.UnmanagedSourceView,
) core.ResultList {
	return list.FromSingle(core.Result{
		{
			Type:     core.ErrorResult,
			Message:  "Explicit register type does not match previous declaration",
			Location: &NewDeclaration,
		},
		{
			Type:     core.HintResult,
			Message:  "Previous declaration here",
			Location: &FirstDeclaration,
		},
	})
}

func NewRegisterAlreadyDefinedResult(
	NewDeclaration core.UnmanagedSourceView,
	FirstDeclaration core.UnmanagedSourceView,
) core.ResultList {
	return list.FromSingle(core.Result{
		{
			Type:     core.ErrorResult,
			Message:  "Register already defined",
			Location: &NewDeclaration,
		},
		{
			Type:     core.HintResult,
			Message:  "Previous definition here",
			Location: &FirstDeclaration,
		},
	})
}
