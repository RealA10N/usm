package tokens

import "fmt"

// currently only supporting i32 immediate(s).
type ImmediateToken struct {
	value int32
}

func (token ImmediateToken) String() string {
	return fmt.Sprintf("<Imm %v>", token.value)
}
