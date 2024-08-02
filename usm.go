package main

import (
	"bytes"
	"fmt"
	"os"
	"usm/lex"
)

func main() {
	file, ok := os.Open("a.usm")
	if ok != nil {
		panic("can't open file")
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		panic("can't read file")
	}

	tokenizer := lex.NewTokenizer()
	view := lex.NewSourceView(buf.String())
	tokens, err := tokenizer.Tokenize(view)
	if err != nil {
		panic(err)
	}

	_, ctx := view.Detach()
	for _, tkn := range tokens {
		fmt.Println(tkn.String(ctx))
	}
}
