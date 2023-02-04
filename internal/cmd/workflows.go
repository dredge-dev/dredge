package cmd

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/dredge-dev/dredge/internal/workflow"
	"github.com/spf13/cobra"
)

func initWorkflowsCommands(de *exec.DredgeExec, rootCmd *cobra.Command) error {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "exec <source>",
		Short: "Execute a remote workflow",
		Long:  "Execute a workflow from a remote Dredgefile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExecCommand(de, args)
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init <source>",
		Short: "Run a remote init workflow",
		Long:  "Execute the init workflow from a remote Dredgefile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInitCommand(de, args)
		},
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

func runInitCommand(e *exec.DredgeExec, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("not enough arguments: missing source")
	}

	de, err := e.Import(config.SourcePath(args[0]))
	if err != nil {
		return err
	}

	w, err := de.GetWorkflow("", "init")
	if err != nil {
		return err
	}

	return workflow.ExecuteWorkflow(w)
}

func runExecCommand(e *exec.DredgeExec, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("not enough arguments: missing source")
	}

	de, err := e.Import(config.SourcePath(args[0]))
	if err != nil {
		return err
	}

	if len(args) == 1 {
		newCmd := &cobra.Command{
			Use: "exec",
		}
		addWorkflows(de, newCmd)
		return newCmd.Help()
	} else {
		var w *exec.Workflow
		if len(args) == 2 {
			w, _ = de.GetWorkflow("", args[1])
			if w == nil {
				b, err := de.GetBucket(args[1])
				if err != nil {
					return err
				}
				cmd, err := createBucketCommand(b)
				if err != nil {
					return err
				}
				return cmd.Help()
			}
			return workflow.ExecuteWorkflow(w)
		}
		w, err = de.GetWorkflow(args[1], args[2])
		if err != nil {
			return err
		}
		return workflow.ExecuteWorkflow(w)
	}
}
