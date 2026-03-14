package lex

import "alon.kr/x/usm/core"

// Comment represents a single-line comment (from ';' to end of line).
// The View covers from the ';' character up to (but not including) the trailing '\n'.
type Comment struct {
	View core.UnmanagedSourceView
}
