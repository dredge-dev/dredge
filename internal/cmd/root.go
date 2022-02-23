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

func Init(dredgeFile *config.DredgeFile) error {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "exec <source>",
		Short: "Execute a remote workflow",
		Long:  "Execute a workflow from a remote Dredgefile",
		RunE:  runExecCommand,
	})
	return addWorkflows(dredgeFile, rootCmd)
}

func addWorkflows(dredgeFile *config.DredgeFile, cmd *cobra.Command) error {
	for _, w := range dredgeFile.Workflows {
		subCmd, err := createWorkflowCommand(dredgeFile, w)
		if err != nil {
			return err
		}
		cmd.AddCommand(subCmd)
	}
	for _, b := range dredgeFile.Buckets {
		subCmd, err := createBucketCommand(dredgeFile, b)
		if err != nil {
			return err
		}
		cmd.AddCommand(subCmd)
	}
	return nil
}

func createWorkflowCommand(dredgeFile *config.DredgeFile, w config.Workflow) (*cobra.Command, error) {
	description, err := dredgeFile.GetWorkflowDescription(w)
	if err != nil {
		return nil, err
	}
	return &cobra.Command{
		Use:   w.Name,
		Short: description,
		Long:  description,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflow.ExecuteWorkflow(dredgeFile, w)
		},
	}, nil
}

func createBucketCommand(dredgeFile *config.DredgeFile, b config.Bucket) (*cobra.Command, error) {
	description, err := dredgeFile.GetBucketDescription(b)
	if err != nil {
		return nil, err
	}
	command := &cobra.Command{
		Use:   b.Name,
		Short: description,
		Long:  description,
	}
	workflows, err := dredgeFile.GetWorkflows(b)
	if err != nil {
		return nil, err
	}
	for _, w := range workflows {
		subCmd, err := createWorkflowCommand(dredgeFile, w)
		if err != nil {
			return nil, err
		}
		command.AddCommand(subCmd)
	}
	return command, nil
}

func Execute() error {
	return rootCmd.Execute()
}
