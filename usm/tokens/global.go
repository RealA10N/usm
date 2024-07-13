package tokens

import "fmt"

type GlobalToken struct {
	name string
}

func (tkn GlobalToken) String() string {
	return fmt.Sprintf("<Global %v>", tkn.name)
}
