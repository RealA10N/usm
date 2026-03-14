package lex

// TokenizeResult holds the output of a Tokenize call: the main token stream
// (with comments stripped) and the captured comments sorted by source position.
type TokenizeResult struct {
	Tokens   []Token
	Comments []Comment
}
