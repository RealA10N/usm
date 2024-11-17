package core

// MARK: UsmError

// The base error interface for all user-facing USM errors.
type UsmError interface {
	// Returns a human-readable, short error message.
	Error(SourceContext) string

	// Returns a hint that can be displayed to the user to help them understand
	// the error. The hint can be empty if there is no hint.
	Hint(SourceContext) string

	// Returns the location of the error in the source code.
	// The view can be of length zero to indicate a single point, or a non-zero
	// view to indicate a range of text.
	Location() UnmanagedSourceView

	// Returns True if the error is internal ("panic"). Internal errors are
	// errors that are caused by a bug in the compiler, and not by the user.
	IsInternalError() bool
}

// MARK: GenericError

// GenericError is a primitive UsmError implementation that uses strings that
// are known in compile time for the error message and the hint.
type GenericError struct {
	ErrorMessage  string
	HintMessage   string
	ErrorLocation UnmanagedSourceView
	IsInternal    bool
}

func (e GenericError) Error(SourceContext) string {
	return e.ErrorMessage
}

func (e GenericError) Hint(SourceContext) string {
	return e.HintMessage
}

func (e GenericError) Location() UnmanagedSourceView {
	return e.ErrorLocation
}

func (e GenericError) IsInternalError() bool {
	return e.IsInternal
}
