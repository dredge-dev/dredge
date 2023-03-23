package cmd

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/api"
	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "drg",
	Short: "Dredge",
	Long:  "Dredge automates DevOps workflows.",
}

func init() {
	rootCmd.AddGroup(&cobra.Group{
		ID:    "resource",
		Title: "Resource Commands:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "workflow",
		Title: "Workflow Commands:",
	})
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print verbose output")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Init(de *exec.DredgeExec) error {
	rootCmd.Long = "Dredge automates DevOps workflows." + GetResourcesHelp(de)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize a Dredge project",
		Long:  "Initialize a Dredge project by creating a Dredgefile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInitCommand(de, args)
		},
	})
	if err := addWorkflowsCommands(de, rootCmd); err != nil {
		return err
	}
	return addResourceCommands(de, rootCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func runInitCommand(de *exec.DredgeExec, args []string) error {
	reset := "\033[0m"
	yellow := "\033[33m"

	// TODO Only add hello workflow if Dredgefile doesn't exist!

	fmt.Printf("Welcome to Dredge!\n")
	fmt.Printf("init creates a Dredgefile, adds a 'hello' workflow, discovers the tools you're using and creates resources for them.\n\n") // TODO Add nice message here -- Should go through DredgeExec and callbacks

	fmt.Printf("\u25B6 Adding the hello workflow... ") // TODO Create start task, update task
	de.AddWorkflowToDredgefile(config.Workflow{
		Name:        "hello",
		Description: "Say hello",
		Steps: []config.Step{
			{
				Shell: &config.ShellStep{
					Cmd: "echo hello",
				},
			},
		},
	})
	err := config.WriteDredgeFile(de.DredgeFile, de.Source)
	if err != nil {
		return err
	}
	fmt.Printf("\u2705\n\n")

	providers, err := de.GetProviders()
	if err != nil {
		return err
	}

	fmt.Printf("\u25B6 Discovering resources...\n")
	for _, provider := range providers {
		err := provider.Discover(de)
		if err != nil {
			de.Log(api.Debug, "Discovery for %s failed with %v", provider.Name(), err)
		}
	}
	fmt.Printf("\u2705\n\n")

	fmt.Printf("\u23E9 Examples to start using Dredge:\n\n")
	fmt.Printf("    %sdrg hello%s                          Executes the hello workflow\n", yellow, reset)
	for resourceName, _ := range de.DredgeFile.Resources {
		if resourceName == "doc" {
			fmt.Printf("    %sdrg search doc --text 'design'%s     Search the docs\n", yellow, reset)
		} else {
			fmt.Printf("    %sdrg get %-27s%sGet the list of %ss\n", yellow, resourceName, reset, resourceName)
		}
	}
	fmt.Printf("\n")

	return nil
}
