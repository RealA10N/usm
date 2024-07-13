package main

import (
	"fmt"
	"os"
	"usm/parser"
)

func main() {
	file, ok := os.Open("a.usm")
	if ok != nil {
		panic("file not found")
	}
	defer file.Close()

	p := parser.Parser{Reader: file}
	tokens, err := p.Parse()
	if err != nil {
		panic(err)
	}

	for _, tkn := range tokens {
		fmt.Println(tkn)
	}
}
