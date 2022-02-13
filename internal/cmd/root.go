package cmd

import (
	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/workflow"
	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "drg",
	Short: "Dredge - toil less, code more",
	Long:  `Dredge automates developer workflows.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print verbose output")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Init(dredgeFile *config.DredgeFile) {
	for _, w := range dredgeFile.Workflows {
		rootCmd.AddCommand(createWorkflowCommand(dredgeFile, w))
	}
	for _, b := range dredgeFile.Buckets {
		command := &cobra.Command{
			Use:   b.Name,
			Short: b.Description,
			Long:  b.Description,
		}
		for _, w := range b.Workflows {
			command.AddCommand(createWorkflowCommand(dredgeFile, w))
		}
		rootCmd.AddCommand(command)
	}
}

func createWorkflowCommand(dredgeFile *config.DredgeFile, w config.Workflow) *cobra.Command {
	return &cobra.Command{
		Use:   w.Name,
		Short: w.Description,
		Long:  w.Description,
		Run: func(cmd *cobra.Command, args []string) {
			err := workflow.ExecuteWorkflow(dredgeFile, w)
			if err != nil {
				panic(err)
			}
		},
	}
}

func Execute() error {
	return rootCmd.Execute()
}
