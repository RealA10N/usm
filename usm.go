package main

import (
	"fmt"
	"os"
	"path/filepath"

	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	usm64core "alon.kr/x/usm/usm64/core"
	"alon.kr/x/usm/usm64/managers"
	usm64ssa "alon.kr/x/usm/usm64/ssa"
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

func printResultAndExit(sourceView core.SourceView, result core.Result) {
	stringer := core.NewResultStringer(sourceView.Ctx(), inputFilepath)
	fmt.Print(stringer.StringResult(result))
	os.Exit(1)
}

func printResultsAndExit(sourceView core.SourceView, results core.ResultList) {
	stringer := core.NewResultStringer(sourceView.Ctx(), inputFilepath)
	for result := range results.Range() {
		fmt.Print(stringer.StringResult(result))
	}
	os.Exit(1)
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
	node, result := parse.NewFileParser().Parse(&tknView)
	if result != nil {
		printResultAndExit(view, result)
	}

	strCtx := parse.StringContext{SourceContext: view.Ctx()}
	fmt.Print(node.String(&strCtx))
}

func emuCommand(cmd *cobra.Command, args []string) {
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
	node, result := parse.NewFileParser().Parse(&tknView)
	if result != nil {
		printResultAndExit(view, result)
	}

	ctx := managers.NewGenerationContext()
	generator := gen.NewFileGenerator()
	info, results := generator.Generate(ctx, view.Ctx(), node)
	if !results.IsEmpty() {
		printResultsAndExit(view, results)
	}

	emulator := usm64core.NewEmulator()
	results = emulator.Emulate(info.Functions[0])
	if !results.IsEmpty() {
		printResultsAndExit(view, results)
	}

	os.Exit(0)
}

func ssaCommand(cmd *cobra.Command, args []string) {
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
	node, result := parse.NewFileParser().Parse(&tknView)
	if result != nil {
		printResultAndExit(view, result)
	}

	ctx := managers.NewGenerationContext()
	generator := gen.NewFileGenerator()
	info, results := generator.Generate(ctx, view.Ctx(), node)
	if !results.IsEmpty() {
		printResultsAndExit(view, results)
	}

	for _, function := range info.Functions {
		results = usm64ssa.ConvertToSsaForm(function)
		if !results.IsEmpty() {
			printResultsAndExit(view, results)
		}
	}

	fmt.Print(info.String())
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

	emuCmd := &cobra.Command{
		Use:     "emu [file]",
		Short:   "Emulate a the main function execution",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: setInputSource,
		Run:     emuCommand,
	}

	ssaCmd := &cobra.Command{
		Use:     "ssa [file]",
		Short:   "Convert a source to SSA form",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: setInputSource,
		Run:     ssaCommand,
	}
	rootCmd.AddCommand(lexCmd)
	rootCmd.AddCommand(fmtCmd)
	rootCmd.AddCommand(emuCmd)
	rootCmd.AddCommand(ssaCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
