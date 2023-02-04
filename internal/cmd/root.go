package cmd

import (
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
	if err := initWorkflowsCommands(de, rootCmd); err != nil {
		return err
	}
	return initResourceCommands(de, rootCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
