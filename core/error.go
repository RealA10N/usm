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
}

// MARK: InternalError

type InternalError struct {
	ErrorMessage  string
	ErrorLocation UnmanagedSourceView
}

func (e InternalError) Error(SourceContext) string {
	return e.ErrorMessage
}

func (InternalError) Hint(SourceContext) string {
	return ""
}

func (e InternalError) Location() UnmanagedSourceView {
	return e.ErrorLocation
}
