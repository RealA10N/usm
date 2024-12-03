package gen

import "alon.kr/x/usm/core"

func NewUndefinedTypeResult(location core.UnmanagedSourceView) core.Result {
	return core.Result{{
		Type:     core.ErrorResult,
		Message:  "Undefined type",
		Location: &location,
	}}
}

func NewUndefinedRegisterResult(location core.UnmanagedSourceView) core.Result {
	return core.Result{{
		Type:     core.ErrorResult,
		Message:  "Undefined register",
		Location: &location,
	}}
}

func NewRegisterTypeMismatchResult(
	NewDeclaration core.UnmanagedSourceView,
	FirstDeclaration core.UnmanagedSourceView,
) core.Result {
	return core.Result{
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
	}
}
