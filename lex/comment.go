package lex

import "alon.kr/x/usm/core"

// Comment represents a single-line comment (from ';' to end of line).
type Comment struct {
	// View covers from the ';' character up to (but not including) the trailing '\n'.
	View core.UnmanagedSourceView
}
