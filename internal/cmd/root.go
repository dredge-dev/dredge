package cmd

import (
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "drg",
	Short: "Dredge",
	Long:  `Dredge automates DevOps workflows.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print verbose output")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Init(de *exec.DredgeExec) error {
	if err := initWorkflowsCommands(de, rootCmd); err != nil {
		return err
	}
	return initResourceCommands(de, rootCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
