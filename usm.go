package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/source"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "usm",
		Short: " One Universal assembly language to rule them all.",
	}

	lexCmd := &cobra.Command{
		Use:   "lex",
		Short: "Lex the source code",
		Run: func(cmd *cobra.Command, args []string) {
			view, err := source.ReadSource(os.Stdin)
			if err != nil {
				panic(err)
			}

			tokens, err := lex.NewTokenizer().Tokenize(view)
			if err != nil {
				panic(err)
			}

			_, ctx := view.Detach()
			for _, tkn := range tokens {
				fmt.Printf("%s ", tkn.String(ctx))
				if tkn.Type == lex.SeparatorToken {
					fmt.Println()
				}
			}
		},
	}

	fmtCmd := &cobra.Command{
		Use:   "fmt",
		Short: "Format the source code",
		Run: func(cmd *cobra.Command, args []string) {
			view, err := source.ReadSource(os.Stdin)
			if err != nil {
				panic(err)
			}

			tokens, err := lex.NewTokenizer().Tokenize(view)
			if err != nil {
				panic(err)
			}

			tknView := parse.NewTokenView(tokens)
			file, perr := parse.NewFileParser().Parse(&tknView)
			if perr == nil {
				fmt.Print(file.String(view.Ctx()))
			} else {
				fmt.Println(perr.Error(view.Ctx()))
			}
		},
	}

	rootCmd.AddCommand(lexCmd)
	rootCmd.AddCommand(fmtCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
