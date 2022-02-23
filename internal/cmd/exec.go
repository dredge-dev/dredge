package cmd

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/workflow"
	"github.com/spf13/cobra"
)

func runExecCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Not enough arguments: missing source")
	}

	dredgeFile, err := config.GetDredgeFile(args[0])
	if err != nil {
		return err
	}

	if len(args) == 1 {
		newCmd := &cobra.Command{
			Use: "exec",
		}
		addWorkflows(dredgeFile, newCmd)
		return newCmd.Help()
	} else {
		var w *config.Workflow
		if len(args) == 2 {
			w, _ = dredgeFile.GetWorkflow("", args[1])
			if w == nil {
				b, _ := dredgeFile.GetBucket(args[1])
				if err != nil {
					return err
				}
				cmd, err := createBucketCommand(dredgeFile, *b)
				if err != nil {
					return err
				}
				return cmd.Help()
			}
			return workflow.ExecuteWorkflow(dredgeFile, *w)
		}
		w, err = dredgeFile.GetWorkflow(args[1], args[2])
		if err != nil {
			return err
		}
		return workflow.ExecuteWorkflow(dredgeFile, *w)
	}
}
