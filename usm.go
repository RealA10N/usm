package main

import (
	"fmt"
	"os"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"
)

func main() {
	file, ok := os.Open("a.usm")
	if ok != nil {
		panic("can't open file")
	}
	defer file.Close()

	tokenizer := lex.NewTokenizer()
	view, err := source.ReadSource(file)
	if err != nil {
		panic("can't read file")
	}

	tokens, err := tokenizer.Tokenize(view)
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Tokens ===")

	_, ctx := view.Detach()
	for _, tkn := range tokens {
		fmt.Printf("%s ", tkn.String(ctx))
		if tkn.Type == lex.SeparatorToken {
			fmt.Println()
		}
	}

	fmt.Println("\n=== Formatted Source ===")

	tknView := parse.NewTokenView(tokens)
	fn, perr := parse.NewFileParser().Parse(&tknView)
	if perr == nil {
		fmt.Print(fn.String(ctx))
	} else {
		fmt.Println(perr.Error(ctx))
	}
}
