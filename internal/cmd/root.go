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

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func Init(dredgeFile *config.DredgeFile) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "exec <target>",
		Short: "Execute a remote workflow",
		Long:  "Execute a workflow from a remote Dredgefile",
		RunE:  runExecCommand,
	})
	addWorkflows(dredgeFile, rootCmd)
}

func addWorkflows(dredgeFile *config.DredgeFile, targetCmd *cobra.Command) {
	for _, w := range dredgeFile.Workflows {
		targetCmd.AddCommand(createWorkflowCommand(dredgeFile, w))
	}
	for _, b := range dredgeFile.Buckets {
		targetCmd.AddCommand(createBucketCommand(dredgeFile, b))
	}
}

func createWorkflowCommand(dredgeFile *config.DredgeFile, w config.Workflow) *cobra.Command {
	return &cobra.Command{
		Use:   w.Name,
		Short: w.Description,
		Long:  w.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflow.ExecuteWorkflow(dredgeFile, w)
		},
	}
}

func createBucketCommand(dredgeFile *config.DredgeFile, b config.Bucket) *cobra.Command {
	command := &cobra.Command{
		Use:   b.Name,
		Short: b.Description,
		Long:  b.Description,
	}
	for _, w := range b.Workflows {
		command.AddCommand(createWorkflowCommand(dredgeFile, w))
	}
	return command
}

func Execute() error {
	return rootCmd.Execute()
}
