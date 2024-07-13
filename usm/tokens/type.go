package tokens

import "fmt"

type TypeToken struct {
	name string
}

func (tkn TypeToken) String() string {
	return fmt.Sprintf("<Type %v>", tkn.name)
}
