package tokens

import "fmt"

type LabelToken struct {
	name string
}

func (tkn LabelToken) String() string {
	return fmt.Sprintf("<Label %v>", tkn.name)
}
