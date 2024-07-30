package base

type WordTokenizer interface {
	Tokenize(word string) (Token, error)
}
