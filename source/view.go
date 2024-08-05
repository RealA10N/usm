package source

import (
	"io"

	"alon.kr/x/view"
)

// Represents a view into a single source file ("substring" of the file source).
//
// This structure is self contained and contains a pointer to the file structure
// and whole data, and is a bit larger in memory than the UnmanagedSourceView.
type SourceView = view.View[rune, uint32]

// Represents a view into a single source file ("substring" of the file source).
//
// Does not store the file content itself, but only the start and end indecies
// of the substring (to not waste memory).
// Use the SourceView type to store a view with context to a specific file.
type UnmanagedSourceView = view.UnmanagedView[rune, uint32]

// The context of a file.
//
// When paired with an UnamangedSourceView, creates a SourceView which
// represents a concrete string from a source file.
type SourceContext = view.ViewContext[rune]

func NewSourceView(data string) SourceView {
	return view.NewView[rune, uint32]([]rune(data))
}

func ReadSource(reader io.Reader) (view SourceView, err error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	view = NewSourceView(string(data))
	return
}
