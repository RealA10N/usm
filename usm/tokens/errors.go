package tokens

import "fmt"

type ErrTokenNotMatched struct {
	Word string
}

func (err ErrTokenNotMatched) Error() string {
	return fmt.Sprintf("word '%s' can't match token type", err.Word)
}

type ErrUnexpectedToken struct{ Word string }

func (err ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token '%s'", err.Word)
}
