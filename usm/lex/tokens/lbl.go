package tokens

import "fmt"

type LblToken struct {
	name string
}

func (token LblToken) String() string {
	return fmt.Sprintf("<Lbl %v>", token.name)
}
