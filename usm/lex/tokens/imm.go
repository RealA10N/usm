package tokens

import "fmt"

// currently only supporting i32 immediate(s).
type ImmToken struct {
	Value int32
}

func (token ImmToken) String() string {
	return fmt.Sprintf("<Imm %v>", token.Value)
}
