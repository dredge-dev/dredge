package cmd

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/dredge-dev/dredge/internal/workflow"
	"github.com/spf13/cobra"
)

func runExecCommand(e *exec.DredgeExec, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Not enough arguments: missing source")
	}

	de, err := e.Import(args[0])
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
