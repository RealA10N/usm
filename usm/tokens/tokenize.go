package tokens

type Tokenizer interface {
	Tokenize(word string) (Token, error)
}
