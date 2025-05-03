package main

import (
	"fmt"
	"os"

	"alon.kr/x/usm/transform"
	"github.com/spf13/cobra"
)

var instructionSets = transform.NewIsaCollection(
	&transform.InstructionSet{
		Name:       "usm",
		Aliases:    nil,
		Extensions: []string{".usm"},
		// GenerationContext: ,
		Transformations: *transform.NewTransformationCollection(
			&transform.Transformation{
				Name:    "dead-code-elimination",
				Aliases: []string{"dce"},
				Target:  "usm",
				// Transform: ,
			},
			&transform.Transformation{
				Name:    "aarch64",
				Aliases: []string{"arm64"},
				Target:  "aarch64",
				// Transform: ,
			},
		),
	},

	&transform.InstructionSet{
		Name:       "aarch64",
		Aliases:    []string{"arm64"},
		Extensions: []string{".aarch64.usm", ".arm64.usm"},
		// GenerationContext: ,
		Transformations: *transform.NewTransformationCollection(
			&transform.Transformation{
				Name:    "dead-code-elimination",
				Aliases: []string{"dce"},
				Target:  "aarch64",
				// Transform: ,
			},
		),
	},
)

func ValidArgsFunction(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return []string{}, cobra.ShellCompDirectiveDefault
	}

	isa := instructionSets.FilenameToInstructionSet(args[0])
	if isa == nil {
		return instructionSets.TransformationNames(), cobra.ShellCompDirectiveNoFileComp
	}

	currentIsa := instructionSets.Traverse(isa, args[1:])
	if currentIsa == nil {
		return instructionSets.TransformationNames(), cobra.ShellCompDirectiveNoFileComp
	}

	return currentIsa.Transformations.Names(), cobra.ShellCompDirectiveNoFileComp
}

func Args(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("requires an input file argument")
	}

	return nil
}

func Run(cmd *cobra.Command, args []string) {
	fmt.Println("Running usm with args:", args)
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
