package cmd

import (
	"github.com/dredge-dev/dredge/internal/exec"
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

func Init(de *exec.DredgeExec) error {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "exec <source>",
		Short: "Execute a remote workflow",
		Long:  "Execute a workflow from a remote Dredgefile",
		RunE:  runExecCommand,
	})
	return addWorkflows(de, rootCmd)
}

func addWorkflows(de *exec.DredgeExec, cmd *cobra.Command) error {
	workflows, err := de.GetWorkflows()
	if err != nil {
		return err
	}
	for _, w := range workflows {
		subCmd, err := createWorkflowCommand(w)
		if err != nil {
			return err
		}
		cmd.AddCommand(subCmd)
	}
	buckets, err := de.GetBuckets()
	if err != nil {
		return err
	}
	for _, b := range buckets {
		subCmd, err := createBucketCommand(b)
		if err != nil {
			return err
		}
		cmd.AddCommand(subCmd)
	}
	return nil
}

func createWorkflowCommand(w *exec.Workflow) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   w.Name,
		Short: w.Description,
		Long:  w.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflow.ExecuteWorkflow(w)
		},
	}, nil
}

func createBucketCommand(b *exec.Bucket) (*cobra.Command, error) {
	command := &cobra.Command{
		Use:   b.Name,
		Short: b.Description,
		Long:  b.Description,
	}
	workflows, err := b.GetWorkflows()
	if err != nil {
		return nil, err
	}
	for _, w := range workflows {
		subCmd, err := createWorkflowCommand(w)
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
