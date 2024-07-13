package tokens

import "fmt"

type RegisterToken struct {
	name string
}

func (tkn RegisterToken) String() string {
	return fmt.Sprintf("<Register %v>", tkn.name)
}
