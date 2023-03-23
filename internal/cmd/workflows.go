package cmd

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/dredge-dev/dredge/internal/workflow"
	"github.com/spf13/cobra"
)

func addWorkflowsCommands(e *exec.DredgeExec, rootCmd *cobra.Command) error {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "exec <source>",
		Short: "Execute a remote workflow",
		Long:  "Execute a workflow from a remote Dredgefile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExecCommand(e, args)
		},
	})
	return addWorkflows(e, rootCmd)
}

func addWorkflows(e *exec.DredgeExec, cmd *cobra.Command) error {
	workflows, err := e.GetWorkflows()
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
	buckets, err := e.GetBuckets()
	if err != nil {
		return err
	}
	for _, b := range buckets {
		subCmd, err := createBucketCommand(e, b)
		if err != nil {
			return err
		}
		cmd.AddCommand(subCmd)
	}
	return nil
}

func createWorkflowCommand(w *workflow.Workflow) (*cobra.Command, error) {
	return &cobra.Command{
		Use:     w.Name,
		Short:   w.Description,
		Long:    w.Description,
		GroupID: "workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.Execute()
		},
	}, nil
}

func createBucketCommand(e *exec.DredgeExec, b *workflow.Bucket) (*cobra.Command, error) {
	command := &cobra.Command{
		Use:     b.Name,
		Short:   b.Description,
		Long:    b.Description,
		GroupID: "workflow",
	}
	command.AddGroup(&cobra.Group{
		ID:    "workflow",
		Title: "Workflows:",
	})
	workflows, err := e.GetWorkflowsInBucket(b)
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
		var w *workflow.Workflow
		if len(args) == 2 {
			w, _ = de.GetWorkflow("", args[1])
			if w == nil {
				b, err := de.GetBucket(args[1])
				if err != nil {
					return err
				}
				cmd, err := createBucketCommand(e, b)
				if err != nil {
					return err
				}
				return cmd.Help()
			}
			return w.Execute()
		}
		w, err = de.GetWorkflow(args[1], args[2])
		if err != nil {
			return err
		}
		return w.Execute()
	}
}
