package main

import (
	"fmt"
	"os"
	"strings"

	aarch64managers "alon.kr/x/usm/aarch64/managers"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/lex"
	"alon.kr/x/usm/parse"
	"alon.kr/x/usm/transform"
	usm64managers "alon.kr/x/usm/usm64/managers"
	"github.com/spf13/cobra"
)

var targets = transform.NewTargetCollection(
	&transform.Target{
		Names:             []string{"usm"},
		Extensions:        []string{".usm"},
		Description:       "A universal assembly language",
		GenerationContext: usm64managers.NewGenerationContext(),
		Transformations: *transform.NewTransformationCollection(
			&transform.Transformation{
				Names:       []string{"dead-code-elimination", "dce"},
				Description: "An optimization pass that eliminates unnecessary instructions",
				TargetName:  "usm",
				// Transform: ,
			},
			&transform.Transformation{
				Names:       []string{"aarch64", "arm64"},
				Description: "Converts the universal assembly to matching machine specific aarch64 assembly",
				TargetName:  "aarch64",
				// Transform: ,
			},
		),
	},

	&transform.Target{
		Names:             []string{"aarch64", "arm64"},
		Extensions:        []string{".aarch64.usm", ".arm64.usm"},
		Description:       "Aarch64 (arm64v8) assembly language",
		GenerationContext: aarch64managers.NewGenerationContext(),
		Transformations: *transform.NewTransformationCollection(
			&transform.Transformation{
				Names:      []string{"macho", "macho-obj", "macho-object"},
				TargetName: "aarch64-macho-object",
				Transform:  aarch64translation.ToMachoObject,
			},
		),
	},

	&transform.Target{
		Names: []string{
			"aarch64-macho-object",
			"aarch64-macho-obj",
			"arm64-macho-object",
			"arm64-macho-obj",
		},
		Extensions:  []string{".o"},
		Description: "Mach-O object file containing aarch64 assembly",
	},
)

var inputFilepath string

func printResultAndExit(sourceView core.SourceView, result core.Result) {
	stringer := core.NewResultStringer(sourceView.Ctx(), inputFilepath)
	fmt.Fprint(os.Stderr, stringer.StringResult(result))
	os.Exit(1)
}

func printResultsAndExit(sourceView core.SourceView, results core.ResultList) {
	stringer := core.NewResultStringer(sourceView.Ctx(), inputFilepath)
	for result := range results.Range() {
		fmt.Fprint(os.Stderr, stringer.StringResult(result))
	}
	os.Exit(1)
}

func ValidArgsFunction(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	// TODO: improve this implementation to ignore flags.

	if len(args) == 0 {
		// Filename is not provided: regular shell completion for filenames.
		return []string{}, cobra.ShellCompDirectiveDefault
	}

	start, _ := targets.FilenameToTarget(args[0])
	if start == nil {
		// Invalid start target name: just suggest all transformation.
		return targets.TransformationNames(), cobra.ShellCompDirectiveNoFileComp
	}

	_, end, results := targets.Traverse(start, args[1:])
	if !results.IsEmpty() {
		// Invalid transformation chain: just suggest all transformations.
		return targets.TransformationNames(), cobra.ShellCompDirectiveNoFileComp
	}

	return end.Transformations.Names(), cobra.ShellCompDirectiveNoFileComp
}

func Args(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("requires an input file argument")
	}

	return nil
}

func Run(cmd *cobra.Command, args []string) {
	inputFilepath = args[0]

	inputTarget, inputExt := targets.FilenameToTarget(inputFilepath)
	if inputTarget == nil {
		fmt.Fprintf(
			os.Stderr,
			"Target type can't be determined from filename: %v\n",
			inputFilepath,
		)
		os.Exit(1)
	}

	ctx := inputTarget.GenerationContext
	if ctx == nil {
		fmt.Fprintf(
			os.Stderr,
			"Target type isn't supported as input: %v\n",
			inputTarget.Names[0],
		)
		os.Exit(1)
	}

	file, err := os.Open(inputFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	view, err := core.ReadSource(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading source: %v\n", err)
		os.Exit(1)
	}

	tokens, err := lex.NewTokenizer().Tokenize(view)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error tokenizing: %v\n", err)
		os.Exit(1)
	}

	tknView := parse.NewTokenView(tokens)
	node, result := parse.NewFileParser().Parse(&tknView)
	if result != nil {
		printResultAndExit(view, result)
	}

	generator := gen.NewFileGenerator()
	info, results := generator.Generate(ctx, view.Ctx(), node)
	if !results.IsEmpty() {
		printResultsAndExit(view, results)
	}

	data := transform.NewTargetData(inputTarget, info)
	transformationNames := args[1:]
	data, results = targets.Transform(data, transformationNames)
	if !results.IsEmpty() {
		printResultsAndExit(view, results)
	}

	// TODO: if the transformation chain is empty, use the same output filepath
	// as the input filepath by default. In general, if the input target is the
	// same as the output target, use the same filepath.

	outputTarget := data.Target
	cleanInputFilepath, _ := strings.CutSuffix(inputFilepath, inputExt)
	outputFilepath := cleanInputFilepath + outputTarget.Extensions[0]

	outputFile, err := os.Create(outputFilepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	if _, err := data.WriteTo(outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:               "usm <input_file> [transformation...]",
		Short:             "One Universal assembly language to rule them all.",
		ValidArgsFunction: ValidArgsFunction,
		Args:              Args,
		Run:               Run,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
