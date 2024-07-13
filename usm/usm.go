package main

import (
	"fmt"
	"os"
	"usm/lex"
)

func main() {
	file, ok := os.Open("a.usm")
	if ok != nil {
		panic("file not found")
	}
	defer file.Close()

	p := lex.Tokenizer{Reader: file}
	tokens, err := p.Tokenize()
	if err != nil {
		panic(err)
	}

	for _, tkn := range tokens {
		fmt.Println(tkn)
	}
}
