package base

import "fmt"

// TODO: store tokens as a struct with (start, end) positions
// as a substring of the source file. This removes duplication of info,
// plus it can help with better error messages (line and column of error).
type Token interface {
	fmt.Stringer
}
