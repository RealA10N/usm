package tokens

import "fmt"

// currently only supporting i32 immediate(s).
type ImmediateToken struct {
	value int32
}

func (tkn ImmediateToken) String() string {
	return fmt.Sprintf("<Immediate %v>", tkn.value)
}
