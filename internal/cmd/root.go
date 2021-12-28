package cmd

import (
	"github.com/spf13/cobra"
	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/workflow"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "drg",
	Short: "Dredge - toil less, code more",
	Long: `Dredge automates developer workflows.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print verbose output")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Init(dredgeFile *config.DredgeFile) {
	for _, w := range dredgeFile.Workflows {
		rootCmd.AddCommand(
			&cobra.Command{
				Use:   w.Name,
				Short: w.Description,
				Long:  w.Description,
				Run: func(cmd *cobra.Command, args []string) {
					err := workflow.ExecuteWorkflow(dredgeFile, cmd, args)
					if err != nil {
						panic(err)
					}
				},
			},
		)
	}
}

func Execute() error {
	return rootCmd.Execute()
}
