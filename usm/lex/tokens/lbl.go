package tokens

import "fmt"

type LblToken struct {
	Name string
}

func (token LblToken) String() string {
	return fmt.Sprintf("<Lbl .%v>", token.Name)
}
