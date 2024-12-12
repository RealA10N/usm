package main

import (
	"fmt"
	"os"
	"path/filepath"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"github.com/spf13/cobra"
)

var inputFilepath string = ""

func setInputSource(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		inputFilepath = filepath.Clean(args[0])
		file, err := os.Open(inputFilepath)
		if err != nil {
			return fmt.Errorf("error opening file: %v", err)
		}
		cmd.SetIn(file)
	}
	return nil
}

func lexCommand(cmd *cobra.Command, args []string) {
	view, err := core.ReadSource(cmd.InOrStdin())
	if err != nil {
		fmt.Printf("Error reading source: %v\n", err)
		os.Exit(1)
	}

	tokens, err := lex.NewTokenizer().Tokenize(view)
	if err != nil {
		fmt.Printf("Error tokenizing: %v\n", err)
		os.Exit(1)
	}

	for _, tkn := range tokens {
		fmt.Printf("%s ", tkn.String(view.Ctx()))
		if tkn.Type == lex.SeparatorToken {
			fmt.Println()
		}
	}
}

func fmtCommand(cmd *cobra.Command, args []string) {
	view, err := core.ReadSource(cmd.InOrStdin())
	if err != nil {
		fmt.Printf("Error reading source: %v\n", err)
		os.Exit(1)
	}

	tokens, err := lex.NewTokenizer().Tokenize(view)
	if err != nil {
		fmt.Printf("Error tokenizing: %v\n", err)
		os.Exit(1)
	}

	tknView := parse.NewTokenView(tokens)
	file, result := parse.NewFileParser().Parse(&tknView)
	if result == nil {
		strCtx := parse.StringContext{SourceContext: view.Ctx()}
		fmt.Print(file.String(&strCtx))
	} else {
		stringer := core.NewResultStringer(view.Ctx(), inputFilepath)
		fmt.Print(stringer.StringResult(result))
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "usm",
		Short: "One Universal assembly language to rule them all.",
	}

	lexCmd := &cobra.Command{
		Use:     "lex [file]",
		Short:   "Lex the source code",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: setInputSource,
		Run:     lexCommand,
	}

	fmtCmd := &cobra.Command{
		Use:     "fmt [file]",
		Short:   "Format the source code",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: setInputSource,
		Run:     fmtCommand,
	}

	rootCmd.AddCommand(lexCmd)
	rootCmd.AddCommand(fmtCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
